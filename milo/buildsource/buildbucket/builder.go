// Copyright 2016 The LUCI Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package buildbucket

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/net/context"

	bb "go.chromium.org/luci/buildbucket"
	"go.chromium.org/luci/buildbucket/proto"
	bbv1 "go.chromium.org/luci/common/api/buildbucket/buildbucket/v1"
	"go.chromium.org/luci/common/clock"
	"go.chromium.org/luci/common/data/strpair"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/logging"
	"go.chromium.org/luci/common/sync/parallel"

	"go.chromium.org/luci/milo/common"
	"go.chromium.org/luci/milo/common/model"
	"go.chromium.org/luci/milo/frontend/ui"
)

// BuilderID represents a buildbucket builder.  We wrap the underlying representation
// since we represent builder IDs slightly differently in Milo vs. Buildbucket.
// I.E. Builders can source from either BuildBot or Buildbucket.
type BuilderID struct {
	// Builder_ID is the buildbucket v2 representation of the builder ID.  Note
	// that the v2 representation uses short bucket names.
	buildbucketpb.Builder_ID
}

func NewBuilderID(v1Bucket, builder string) (bid BuilderID) {
	bid.Project, bid.Bucket = bb.BucketNameToV2(v1Bucket)
	bid.Builder = builder
	return
}

// V1Bucket returns the buildbucket v1 representation of the bucket name, which
// is what we use in Milo.
func (b BuilderID) V1Bucket() string {
	return fmt.Sprintf("luci.%s.%s", b.Project, b.Bucket)
}

// String returns the canonical format of BuilderID.
func (b BuilderID) String() string {
	return fmt.Sprintf("buildbucket/%s/%s", b.V1Bucket(), b.Builder)
}

// fetchBuilds fetches builds given a criteria.
// The returned builds are sorted by build creation descending.
// count defines maximum number of builds to fetch; if <0, defaults to 100.
func fetchBuilds(c context.Context, client *bbv1.Service, bid BuilderID,
	status string, limit int) ([]*bbv1.ApiCommonBuildMessage, error) {

	search := client.Search()
	search.Bucket(bid.V1Bucket())
	search.Status(status)
	search.Tag(strpair.Format(bbv1.TagBuilder, bid.Builder))
	search.IncludeExperimental(true)

	if limit < 0 {
		limit = 100
	}

	start := clock.Now(c)
	msgs, err := search.Fetch(limit, nil)
	if err != nil {
		return nil, err
	}
	logging.Infof(c, "Fetched %d %s builds in %s", len(msgs), status, clock.Since(c, start))
	return msgs, nil
}

func getDebugBuilds(c context.Context, bid BuilderID, maxCompletedBuilds int, target *ui.Builder) error {
	// ../buildbucket below assumes that
	// - this code is not executed by tests outside of this dir
	// - this dir is a sibling of frontend dir
	resFile, err := os.Open(filepath.Join(
		"..", "buildbucket", "testdata", bid.V1Bucket(), bid.Builder+".json"))
	if err != nil {
		return err
	}
	defer resFile.Close()

	res := &bbv1.ApiSearchResponseMessage{}
	if err := json.NewDecoder(resFile).Decode(res); err != nil {
		return err
	}

	for _, bb := range res.Builds {
		mb, err := toMiloBuild(c, bb)
		if err != nil {
			return err
		}
		bs := mb.BuildSummary()
		switch mb.Summary.Status {
		case model.NotRun:
			target.PendingBuilds = append(target.PendingBuilds, bs)

		case model.Running:
			target.CurrentBuilds = append(target.CurrentBuilds, bs)

		case model.Success, model.Failure, model.InfraFailure, model.Warning:
			if len(target.FinishedBuilds) < maxCompletedBuilds {
				target.FinishedBuilds = append(target.FinishedBuilds, bs)
			}

		default:
			panic("impossible")
		}
	}
	return nil
}

func getHost(c context.Context) (string, error) {
	settings := common.GetSettings(c)
	if settings.Buildbucket == nil || settings.Buildbucket.Host == "" {
		return "", errors.New("missing buildbucket host in settings")
	}
	return settings.Buildbucket.Host, nil
}

// GetBuilder is used by buildsource.BuilderID.Get to obtain the resp.Builder.
func GetBuilder(c context.Context, bid BuilderID, limit int) (*ui.Builder, error) {
	host, err := getHost(c)
	if err != nil {
		return nil, err
	}

	if limit < 0 {
		limit = 20
	}

	var lock sync.Mutex
	result := &ui.Builder{
		Name: bid.Builder,
	}
	if host == "debug" {
		return result, getDebugBuilds(c, bid, limit, result)
	}
	client, err := newBuildbucketClient(c, host)
	if err != nil {
		return nil, err
	}

	fetch := func(statusFilter string, limit int) error {
		msgs, err := fetchBuilds(c, client, bid, statusFilter, limit)
		if err != nil {
			logging.Errorf(c, "Could not fetch %s builds: %s", statusFilter, err)
			return err
		}
		for _, m := range msgs {
			mb, err := toMiloBuild(c, m)
			if err != nil {
				return errors.Annotate(err, "failed to convert build %d to milo build", m.Id).Err()
			}
			bs := mb.BuildSummary()
			lock.Lock()
			switch mb.Summary.Status {
			case model.NotRun:
				result.PendingBuilds = append(result.PendingBuilds, bs)
			case model.Running:
				result.CurrentBuilds = append(result.CurrentBuilds, bs)
			default:
				result.FinishedBuilds = append(result.FinishedBuilds, bs)
			}
			lock.Unlock()
		}
		return nil
	}
	return result, parallel.FanOutIn(func(work chan<- func() error) {
		work <- func() (err error) {
			result.MachinePool, err = getPool(c, bid)
			return
		}
		work <- func() error {
			return fetch(bbv1.StatusScheduled, -1)
		}
		work <- func() error {
			return fetch(bbv1.StatusStarted, -1)
		}
		work <- func() error {
			return fetch(bbv1.StatusCompleted, limit)
		}
	})
}
