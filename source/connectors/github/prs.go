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

package github

import (
	"context"
	"sync"
	"time"

	"github.com/artem-sidorenko/chagen/source/connectors/helpers"

	"github.com/artem-sidorenko/chagen/data"

	"github.com/google/go-github/github"
)

// PRsPerPage defined how many PRs are fetched per page
var PRsPerPage = 30 // nolint: gochecknoglobals

const (
	prsProcessingRoutines = 10
)

// MRs returns the PRs via channels.
// Returns possible errors via given cerr channel
// cmrs returns MRs
// cmrscounter returns the channel, which ticks when a PR is proceeded
// cmaxmrs returns the max available amount of MRs
func (c *Connector) MRs(
	ctx context.Context,
	cerr chan<- error,
) (
	cmrs <-chan data.MR,
	cmrscounter <-chan bool,
	cmaxmrs <-chan int,
) {
	// for detailed comments, please see the tags.go
	mrs := make(chan []*github.PullRequest)
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

		resp, n, err := c.processPRPage(sctx, 1, mrs)
		if err != nil {
			helpers.NonBlockingErrSend(sctx, scerr, err)
			closeCh()
			return
		}

		if resp.LastPage == 0 {
			maxmrs <- n
		} else {
			cpages := c.processPRPages(sctx, scerr, maxmrs, mrs, &wg, resp.LastPage)

			for i := resp.LastPage; i >= 2; i-- {
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

	dmrs := c.processPRs(ctx, cerr, mrs, mrscounter, &wgT)

	go func() {
		wgTP.Wait()
		wgT.Wait()
		close(scerr)
		close(mrscounter)
	}()

	return dmrs, mrscounter, maxmrs
}

// processPRPage gets the PRs from GitHub for given page and returns them via
// given channel. PRsCount contains the amount of PRs in the current response
func (c *Connector) processPRPage(
	ctx context.Context,
	page int,
	ret chan<- []*github.PullRequest,
) (
	resp *github.Response,
	prsCount int,
	err error,
) {

	prs, resp, err := c.client.PullRequests.List(
		ctx,
		c.Owner,
		c.Repo,
		&github.PullRequestListOptions{
			State:       "closed",
			ListOptions: github.ListOptions{Page: page, PerPage: PRsPerPage},
		},
	)
	if err != nil {
		return nil, 0, err
	}

	select {
	case <-ctx.Done():
		return nil, 0, nil
	case ret <- prs:
		return resp, len(prs), nil
	}
}

// processPRPages processes GitHub PR page numbers, given in the cpages channel and returns
// the GH PullRequest data structures via channel
// possible errors are returned via given cerr channel
func (c *Connector) processPRPages(
	ctx context.Context,
	cerr chan<- error,
	cmaxprs chan<- int,
	prs chan<- []*github.PullRequest,
	wg *sync.WaitGroup,
	lastPage int,
) (cpages chan<- int) {
	ret := make(chan int)

	for i := 0; i < prsProcessingRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for page := range ret {
				_, n, err := c.processPRPage(ctx, page, prs)

				if err != nil {
					helpers.NonBlockingErrSend(ctx, cerr, err)
					return
				}

				if page == lastPage {
					cmaxprs <- n + (lastPage-1)*PRsPerPage
				}
			}
		}()
	}

	return ret
}

// processPRs processes given GitHub PRs in the cprs channel and returns
// the PRs in our data structure via channel
// possible errors are returned via given cerr channel
func (c *Connector) processPRs(
	ctx context.Context,
	_ chan<- error,
	cprs <-chan []*github.PullRequest,
	cmrscounter chan<- bool,
	wg *sync.WaitGroup,
) <-chan data.MR {

	ret := make(chan data.MR)

	for i := 0; i < prsProcessingRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for prs := range cprs {
				for _, pr := range prs {
					cmrscounter <- true
					// we need only merged PRs, skip everything else
					if pr.GetMergedAt() == (time.Time{}) {
						continue
					}

					var lbs []string
					if pr.Labels != nil && len(pr.Labels) > 0 {
						for _, l := range pr.Labels {
							lbs = append(lbs, *l.Name)
						}
					}

					pr := data.MR{
						ID:         pr.GetNumber(),
						Name:       pr.GetTitle(),
						MergedDate: pr.GetMergedAt(),
						URL:        pr.GetHTMLURL(),
						Author:     pr.User.GetLogin(),
						AuthorURL:  pr.User.GetHTMLURL(),
						Labels:     lbs,
					}

					select {
					case <-ctx.Done():
						return
					case ret <- pr:
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
