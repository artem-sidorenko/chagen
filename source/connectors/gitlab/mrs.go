/*
   Copyright 2019 Artem Sidorenko <artem@posteo.de>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package gitlab

import (
	"context"
	"sync"

	gitlab "github.com/xanzy/go-gitlab"

	"github.com/artem-sidorenko/chagen/source/connectors/helpers"

	"github.com/artem-sidorenko/chagen/data"
)

// MRsPerPage defined how many MRs are fetched per page
var MRsPerPage = 30 // nolint: gochecknoglobals

const (
	mrsProcessingRoutines = 10
)

// MRs returns the MRs via channels.
// Returns possible errors via given cerr channel
// cmrs returns MRs
// cmrscounter returns the channel, which ticks when a MR is proceeded
// cmaxmrs returns the max available amount of MRs
func (c *Connector) MRs(
	ctx context.Context,
	cerr chan<- error,
) (
	cmrs <-chan data.MR,
	cmrscounter <-chan bool,
	cmaxmrs <-chan int,
) {
	// for detailed comments, please see the github/tags.go
	mrs := make(chan []*gitlab.MergeRequest)
	maxmrs := make(chan int)
	mrscounter := make(chan bool, 100)

	sctx, cancel := context.WithCancel(ctx)
	scerr := make(chan error)

	var wgTP, wgT sync.WaitGroup

	go func() {
		select {
		case <-ctx.Done():
			return
		case err, ok := <-scerr:
			if ok {
				cancel()
				cerr <- err
			}
		}
	}()

	wgTP.Add(1)
	go func() {
		var wg sync.WaitGroup

		closeCh := func() {
			close(mrs)
			close(maxmrs)
		}

		resp, n, err := c.processMRPage(sctx, 1, mrs)
		if err != nil {
			helpers.NonBlockingErrSend(sctx, scerr, err)
			closeCh()
			return
		}

		if resp.TotalPages == 1 {
			maxmrs <- n
		} else {
			cpages := c.processMRPages(sctx, scerr, maxmrs, mrs, &wg, resp.TotalPages)

			for i := resp.TotalPages; i >= 2; i-- {
				cpages <- i
			}
			close(cpages)
		}

		go func() {
			wg.Wait()
			closeCh()
			wgTP.Done()
		}()
	}()

	dmrs := c.processMRs(ctx, cerr, mrs, mrscounter, &wgT)

	go func() {
		wgTP.Wait()
		wgT.Wait()
		close(scerr)
		close(mrscounter)
	}()

	return dmrs, mrscounter, maxmrs
}

// processMRPage gets the MRs from GitLab for given page and returns them via
// given channel. MRsCount contains the amount of MRs in the current response
func (c *Connector) processMRPage(
	ctx context.Context,
	page int,
	ret chan<- []*gitlab.MergeRequest,
) (
	resp *gitlab.Response,
	mrsCount int,
	err error,
) {

	mrs, resp, err := c.client.MergeRequests.ListProjectMergeRequests(
		c.Owner+"/"+c.Repo,
		&gitlab.ListProjectMergeRequestsOptions{
			State:       helpers.StringPtr("merged"),
			ListOptions: gitlab.ListOptions{Page: page, PerPage: MRsPerPage}},
	)

	if err != nil {
		return nil, 0, err
	}

	select {
	case <-ctx.Done():
		return nil, 0, nil
	case ret <- mrs:
		return resp, len(mrs), nil
	}
}

// processMRPages processes GitHub PR page numbers, given in the cpages channel and returns
// the GH PullRequest data structures via channel
// possible errors are returned via given cerr channel
func (c *Connector) processMRPages(
	ctx context.Context,
	cerr chan<- error,
	cmaxmrs chan<- int,
	mrs chan<- []*gitlab.MergeRequest,
	wg *sync.WaitGroup,
	lastPage int,
) (cpages chan<- int) {
	ret := make(chan int)

	for i := 0; i < mrsProcessingRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for page := range ret {
				_, n, err := c.processMRPage(ctx, page, mrs)

				if err != nil {
					helpers.NonBlockingErrSend(ctx, cerr, err)
					return
				}

				if page == lastPage {
					cmaxmrs <- n + (lastPage-1)*MRsPerPage
				}
			}
		}()
	}

	return ret
}

// processMRs processes given GitLab MRs in the cmrs channel and returns
// the MRs in our data structure via channel
// possible errors are returned via given cerr channel
func (c *Connector) processMRs(
	ctx context.Context,
	cerr chan<- error,
	cmrs <-chan []*gitlab.MergeRequest,
	cmrscounter chan<- bool,
	wg *sync.WaitGroup,
) <-chan data.MR {

	ret := make(chan data.MR)

	for i := 0; i < mrsProcessingRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for mrs := range cmrs {
				for _, mr := range mrs {
					cmrscounter <- true

					authorURL, err := c.getUsernameURL(mr.Author.Username)
					if err != nil {
						helpers.NonBlockingErrSend(ctx, cerr, err)
						return
					}

					rmr := data.MR{
						ID:         mr.IID,
						Name:       mr.Title,
						URL:        mr.WebURL,
						MergedDate: *mr.MergedAt,
						Author:     mr.Author.Username,
						AuthorURL:  authorURL,
						Labels:     mr.Labels,
					}

					select {
					case <-ctx.Done():
						return
					case ret <- rmr:
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ret)
	}()

	return ret
}
