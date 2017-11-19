// Copyright 2017 The LUCI Authors.
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

package model

import (
	"fmt"
	"testing"
	"time"

	"golang.org/x/net/context"

	"go.chromium.org/gae/service/datastore"
	"go.chromium.org/luci/appengine/gaetesting"
	"go.chromium.org/luci/common/clock/testclock"
	"go.chromium.org/luci/server/caching"

	"go.chromium.org/luci/milo/common"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUpdateBuilder(t *testing.T) {
	t.Parallel()

	Convey(`TestUpdateBuilder`, t, func() {
		c := gaetesting.TestingContextWithAppID("luci-milo-dev")

		builder := &BuilderSummary{BuilderID: "fake"}

		// Populate a few BuildSummaries. For convenience, ordered by creation time.
		builds := make([]*BuildSummary, 10)
		for i := 0; i < 10; i++ {
			builds[i] = &BuildSummary{
				BuildKey:  datastore.MakeKey(c, "fakeBuild", i),
				BuilderID: builder.BuilderID,
				BuildID:   fmt.Sprintf("build_id/%d", i),
				Created:   testclock.TestRecentTimeUTC.Add(time.Duration(i) * time.Hour),
			}
		}

		c = caching.WithRequestCache(c)

		updateBuilder := func(build *BuildSummary) {
			err := datastore.RunInTransaction(c, func(c context.Context) error {
				return UpdateBuilderForBuild(c, build)
			}, nil)
			So(err, ShouldBeNil)
			err = datastore.Get(c, builder)
			So(err, ShouldBeNil)
		}

		Convey("Updating appropriate builder having existing last finished build", func() {
			builder.LastFinishedCreated = builds[5].Created
			builder.LastFinishedStatus = Success
			builder.LastFinishedBuildID = builds[5].BuildID
			err := datastore.Put(c, builder)
			So(err, ShouldBeNil)

			Convey("with finished build should not update last finished build info", func() {
				builds[6].Summary.Status = Failure
				updateBuilder(builds[6])
				So(builder.LastFinishedStatus, ShouldEqual, Failure)
				So(builder.LastFinishedBuildID, ShouldEqual, builds[6].BuildID)
			})

			Convey("for build created earlier than last finished", func() {
				builds[4].Summary.Status = Failure
				updateBuilder(builds[4])
				So(builder.LastFinishedStatus, ShouldEqual, Success)
				So(builder.LastFinishedBuildID, ShouldEqual, builds[5].BuildID)
			})

			Convey("for build created later than last finished", func() {
				builds[6].Summary.Status = NotRun
				updateBuilder(builds[6])
				So(builder.LastFinishedStatus, ShouldEqual, Success)
				So(builder.LastFinishedBuildID, ShouldEqual, builds[5].BuildID)
			})
		})

		Convey("Updating appropriate builder with no last finished build should initialize it", func() {
			builds[5].Summary.Status = Failure
			updateBuilder(builds[5])
			So(builder.LastFinishedStatus, ShouldEqual, Failure)
			So(builder.LastFinishedBuildID, ShouldEqual, builds[5].BuildID)
		})
	})
}

func TestBuildIDLink(t *testing.T) {
	t.Parallel()

	Convey(`TestLastFinishedBuildIDLink`, t, func() {
		Convey("Buildbot build gets expected link", func() {
			Convey("with valid BuildID", func() {
				buildID, project := "buildbot/master/builder/number", "proj"
				So(buildIDLink(buildID, project), ShouldEqual, "/buildbot/master/builder/number")
			})

			Convey("with invalid BuildID", func() {
				Convey("with too few tokens", func() {
					buildID, project := "buildbot/wat", "proj"
					So(buildIDLink(buildID, project), ShouldEqual, "#invalid-build-id")
				})

				Convey("with too many tokens", func() {
					buildID, project := "buildbot/wat/wat/wat/wat", "proj"
					So(buildIDLink(buildID, project), ShouldEqual, "#invalid-build-id")
				})
			})
		})

		Convey("Buildbucket build gets expected link", func() {
			Convey("with bucket info", func() {
				buildID, project := "buildbucket/luci.proj.bucket/builder/123", ""
				So(
					buildIDLink(buildID, project),
					ShouldEqual,
					"/p/proj/builders/luci.proj.bucket/builder/123")
			})

			Convey("with only ID info", func() {
				buildID, project := "buildbucket/123", "proj"
				So(buildIDLink(buildID, project), ShouldEqual, "/p/proj/builds/b123")
			})

			Convey("with invalid BuildID", func() {
				Convey("due to missing bucket info", func() {
					buildID, project := "buildbucket/", "proj"
					So(buildIDLink(buildID, project), ShouldEqual, "#invalid-build-id")
				})

				Convey("due to missing project info", func() {
					buildID, project := "buildbucket/123", ""
					So(buildIDLink(buildID, project), ShouldEqual, "#invalid-build-id")
				})
			})
		})

		Convey("Invalid BuildID gets expected link", func() {
			Convey("with unknown source gets expected link", func() {
				buildID, project := "unknown/1", "proj"
				So(buildIDLink(buildID, project), ShouldEqual, "#invalid-build-id")
			})

			Convey("with too few tokens", func() {
				buildID, project := "source", "proj"
				So(buildIDLink(buildID, project), ShouldEqual, "#invalid-build-id")
			})
		})
	})
}

func TestGetBuilderHistories(t *testing.T) {
	t.Parallel()

	Convey(`TestGetBuilderHistories`, t, func() {
		c := gaetesting.TestingContextWithAppID("luci-milo-dev")
		c = caching.WithRequestCache(c)

		datastore.GetTestable(c).AddIndexes(&datastore.IndexDefinition{
			Kind: "BuildSummary",
			SortBy: []datastore.IndexColumn{
				{Property: "BuilderID"},
				{Property: "Created", Descending: true},
			},
		})
		datastore.GetTestable(c).AddIndexes(&datastore.IndexDefinition{
			Kind: "BuildSummary",
			SortBy: []datastore.IndexColumn{
				{Property: "BuilderID"},
				{Property: "Summary.Status"},
			},
		})
		datastore.GetTestable(c).CatchupIndexes()
		datastore.GetTestable(c).Consistent(true)

		nBuilds := 9
		statuses := []Summary{
			{Status: Running},      /* b2 */
			{Status: Success},      /* b2 */
			{Status: Running},      /* b2 */
			{Status: Exception},    /* b2 */
			{Status: Running},      /* b2 */
			{Status: InfraFailure}, /* b2 */
			{Status: NotRun},       /* b2 */
			{Status: NotRun},       /* b2 */
			{Status: Success},      /* b3 */
		}
		So(statuses, ShouldHaveLength, nBuilds)

		p := "proj"
		proj := datastore.MakeKey(c, "Project", p)

		// Populate console.
		err := datastore.Put(c, &common.Console{
			Parent:   proj,
			ID:       "view",
			Builders: []string{"buildbot/master/b1", "buildbot/master/b2"},
		})
		So(err, ShouldBeNil)

		// When populating BuildKey, 0 is invalid, so be sure to increment.
		i := 0

		// Populate builds.
		builds := make([]*BuildSummary, nBuilds)

		// One builder has lots of builds.
		builder := "buildbot/master/b2"
		for i < 8 {
			builds[i] = &BuildSummary{
				BuildKey:  datastore.MakeKey(c, "build", i+1),
				ProjectID: p,
				BuilderID: builder,
				BuildID:   fmt.Sprintf("%s/%d", builder, i),
				Created:   testclock.TestRecentTimeUTC.Add(time.Duration(i) * time.Hour),
				Summary:   statuses[i],
			}
			i++
		}

		// One builder is on a view but has no builds.

		// One builder is not on any project's views.
		builds[i] = &BuildSummary{
			BuildKey:  datastore.MakeKey(c, "build", i+1),
			ProjectID: p,
			BuilderID: "buildbot/master/b3",
			BuildID:   fmt.Sprintf("buildbot/master/b3/%d", i),
			Created:   testclock.TestRecentTimeUTC.Add(time.Duration(i) * time.Hour),
			Summary:   statuses[i],
		}

		err = datastore.Put(c, builds)
		So(err, ShouldBeNil)

		Convey("Getting recent history for existing project", func() {
			Convey("across all views", func() {
				Convey("with limit less than number of finished builds works", func() {
					hists, err := GetBuilderHistories(c, p, 2)
					So(err, ShouldBeNil)
					So(hists, ShouldHaveLength, 2)

					So(*hists[0], ShouldResemble, BuilderHistory{
						BuilderID:    "buildbot/master/b1",
						BuilderLink:  "/buildbot/master/b1",
						RecentBuilds: []*BuildSummary{},
					})

					So(hists[1].BuilderID, ShouldEqual, "buildbot/master/b2")
					So(hists[1].BuilderLink, ShouldEqual, "/buildbot/master/b2")
					So(hists[1].NumPending, ShouldEqual, 2)
					So(hists[1].NumRunning, ShouldEqual, 3)
					So(hists[1].RecentBuilds, ShouldHaveLength, 2)
					So(hists[1].RecentBuilds[0].BuildID, ShouldEqual, builds[5].BuildID)
					So(hists[1].RecentBuilds[1].BuildID, ShouldEqual, builds[3].BuildID)
				})

				Convey("with limit greater than number of finished builds works", func() {
					hists, err := GetBuilderHistories(c, p, 5)
					So(err, ShouldBeNil)
					So(hists, ShouldHaveLength, 2)

					So(*hists[0], ShouldResemble, BuilderHistory{
						BuilderID:    "buildbot/master/b1",
						BuilderLink:  "/buildbot/master/b1",
						RecentBuilds: []*BuildSummary{},
					})

					So(hists[1].BuilderID, ShouldEqual, "buildbot/master/b2")
					So(hists[1].BuilderLink, ShouldEqual, "/buildbot/master/b2")
					So(hists[1].NumPending, ShouldEqual, 2)
					So(hists[1].NumRunning, ShouldEqual, 3)
					So(hists[1].RecentBuilds, ShouldHaveLength, 3)
					So(hists[1].RecentBuilds[0].BuildID, ShouldEqual, builds[5].BuildID)
					So(hists[1].RecentBuilds[1].BuildID, ShouldEqual, builds[3].BuildID)
					So(hists[1].RecentBuilds[2].BuildID, ShouldEqual, builds[1].BuildID)
				})
			})
		})

		Convey("Getting recent history for nonexisting project", func() {
			hists, err := GetBuilderHistories(c, "no_proj", 3)
			So(err, ShouldBeNil)
			So(hists, ShouldHaveLength, 0)
		})
	})
}
