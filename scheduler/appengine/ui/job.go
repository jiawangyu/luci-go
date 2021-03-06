// Copyright 2015 The LUCI Authors.
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

package ui

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"

	mc "go.chromium.org/gae/service/memcache"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/server/router"
	"go.chromium.org/luci/server/templates"

	"go.chromium.org/luci/scheduler/appengine/engine"
	"go.chromium.org/luci/scheduler/appengine/internal"
)

func jobPage(ctx *router.Context) {
	c, w, r, p := ctx.Context, ctx.Writer, ctx.Request, ctx.Params
	e := config(c).Engine

	projectID := p.ByName("ProjectID")
	jobName := p.ByName("JobName")
	cursor := r.URL.Query().Get("c")

	job := jobFromEngine(ctx, projectID, jobName)
	if job == nil {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	// We show active invocations only on the first page of the results.
	var invsActive []*engine.Invocation
	var invsActiveErr error
	var haveInvsActive bool
	if cursor == "" {
		haveInvsActive = true
		wg.Add(1)
		go func() {
			defer wg.Done()
			invsActive, _, invsActiveErr = e.ListInvocations(c, job, engine.ListInvocationsOpts{
				PageSize:   100000000, // ~ ∞, UI doesn't paginate active invocations
				ActiveOnly: true,
			})
		}()
	}

	var invsLog []*engine.Invocation
	var nextCursor string
	var invsLogErr error
	go func() {
		defer wg.Done()
		invsLog, nextCursor, invsLogErr = e.ListInvocations(c, job, engine.ListInvocationsOpts{
			PageSize:     50,
			Cursor:       cursor,
			FinishedOnly: true,
		})
	}()

	var triggers []*internal.Trigger
	var triErr error
	go func() {
		defer wg.Done()
		triggers, triErr = e.ListTriggers(c, job)
	}()

	wg.Wait()

	switch {
	case invsActiveErr != nil:
		panic(errors.Annotate(invsActiveErr, "failed to fetch active invocations").Err())
	case invsLogErr != nil:
		panic(errors.Annotate(invsLogErr, "failed to fetch invocation log").Err())
	case triErr != nil:
		panic(errors.Annotate(triErr, "failed to fetch triggers").Err())
	}

	// memcacheKey hashes cursor to reduce its length, since full cursor doesn't
	// fit into memcache key length limits. Use 'v2' scheme for this ('v1' was
	// used before hashing was added).
	memcacheKey := func(cursor string) string {
		blob := sha1.Sum([]byte(job.JobID + ":" + cursor))
		encoded := base64.StdEncoding.EncodeToString(blob[:])
		return "v2:cursors:list_invocations:" + encoded
	}

	// Cheesy way of implementing bidirectional pagination with forward-only
	// datastore cursors: store mapping from a page cursor to a previous page
	// cursor in the memcache.
	prevCursor := ""
	if cursor != "" {
		if itm, err := mc.GetKey(c, memcacheKey(cursor)); err == nil {
			prevCursor = string(itm.Value())
		}
	}
	if nextCursor != "" {
		itm := mc.NewItem(c, memcacheKey(nextCursor))
		if cursor == "" {
			itm.SetValue([]byte("NULL"))
		} else {
			itm.SetValue([]byte(cursor))
		}
		itm.SetExpiration(24 * time.Hour)
		mc.Set(c, itm)
	}

	// List of invocations in job.ActiveInvocations may contain recently finished
	// invocations not yet removed from the active list by the triage procedure.
	// 'invsActive' is always more accurate, since it fetches invocations from
	// the datastore and checks their status. So update the job entity to be more
	// accurate if we can. This is important for reporting jobs with recently
	// finished invocations as not running. Otherwise the UI page may appear
	// non-consistent (no running invocations in the list, yet the job's status is
	// displayed as "Running"). This is a bit of a hack...
	if haveInvsActive {
		ids := make([]int64, len(invsActive))
		for i, inv := range invsActive {
			ids[i] = inv.ID
		}
		job.ActiveInvocations = ids
	}

	jobUI := makeJob(c, job)
	invsActiveUI := make([]*invocation, len(invsActive))
	for i, inv := range invsActive {
		invsActiveUI[i] = makeInvocation(jobUI, inv)
	}
	invsLogUI := make([]*invocation, len(invsLog))
	for i, inv := range invsLog {
		invsLogUI[i] = makeInvocation(jobUI, inv)
	}

	templates.MustRender(c, w, "pages/job.html", map[string]interface{}{
		"Job":               jobUI,
		"InvocationsActive": invsActiveUI,
		"InvocationsLog":    invsLogUI,
		"PendingTriggers":   makeTriggerList(jobUI.now, triggers),
		"PrevCursor":        prevCursor,
		"NextCursor":        nextCursor,
	})
}

////////////////////////////////////////////////////////////////////////////////
// Actions.

func runJobAction(ctx *router.Context) {
	c, w, r, p := ctx.Context, ctx.Writer, ctx.Request, ctx.Params
	e := config(c).Engine

	projectID := p.ByName("ProjectID")
	jobName := p.ByName("JobName")

	job := jobFromEngine(ctx, projectID, jobName)
	if job == nil {
		return
	}

	// genericReply renders "we did something (or we failed to do something)"
	// page, shown on error or if invocation is starting for too long.
	genericReply := func(err error) {
		templates.MustRender(c, w, "pages/run_job_result.html", map[string]interface{}{
			"ProjectID": projectID,
			"JobName":   jobName,
			"Error":     err,
		})
	}

	// Start the new invocation and redirect the user to its page.
	inv, err := e.ForceInvocation(c, job)
	switch {
	case err == engine.ErrNoPermission:
		uiErrActionForbidden.render(ctx)
	case err != nil:
		genericReply(err)
	default:
		http.Redirect(w, r, fmt.Sprintf("/jobs/%s/%s/%d", projectID, jobName, inv.ID), http.StatusFound)
	}
}

func pauseJobAction(c *router.Context) {
	handleJobAction(c, func(job *engine.Job) error {
		return config(c.Context).Engine.PauseJob(c.Context, job)
	})
}

func resumeJobAction(c *router.Context) {
	handleJobAction(c, func(job *engine.Job) error {
		return config(c.Context).Engine.ResumeJob(c.Context, job)
	})
}

func abortJobAction(c *router.Context) {
	handleJobAction(c, func(job *engine.Job) error {
		return config(c.Context).Engine.AbortJob(c.Context, job)
	})
}

func handleJobAction(c *router.Context, cb func(*engine.Job) error) {
	projectID := c.Params.ByName("ProjectID")
	jobName := c.Params.ByName("JobName")

	job := jobFromEngine(c, projectID, jobName)
	if job == nil {
		return
	}

	switch err := cb(job); {
	case err == engine.ErrNoPermission:
		uiErrActionForbidden.render(c)
	case err != nil:
		panic(err)
	default:
		http.Redirect(c.Writer, c.Request, fmt.Sprintf("/jobs/%s/%s", projectID, jobName), http.StatusFound)
	}
}
