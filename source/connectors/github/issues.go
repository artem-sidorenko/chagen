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

	"github.com/artem-sidorenko/chagen/source/connectors/helpers"

	"github.com/artem-sidorenko/chagen/data"

	"github.com/google/go-github/github"
)

// IssuesPerPage defined how many Issues are fetched per page
var IssuesPerPage = 30 // nolint: gochecknoglobals

const (
	issuesProcessingRoutines = 10
)

// Issues returns the issues via channels.
// Returns possible errors via given cerr channel
// cissues returns issues
// cissuescounter returns the channel, which ticks when an issue is proceeded
// cmaxissues returns the max available amount of issues
func (c *Connector) Issues(
	ctx context.Context,
	cerr chan<- error,
) (
	ctags <-chan data.Issue,
	cissuescounter <-chan bool,
	cmaxissues <-chan int,
) {
	// for detailed comments, please see the tags.go
	issues := make(chan []*github.Issue)
	maxissues := make(chan int)
	issuescounter := make(chan bool, 100)

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
			close(issues)
			close(maxissues)
		}

		resp, n, err := c.processIssuesPage(sctx, 1, issues)
		if err != nil {
			helpers.NonBlockingErrSend(sctx, scerr, err)
			closeCh()
			return
		}

		if resp.LastPage == 0 {
			maxissues <- n
		} else {
			cpages := c.processIssuesPages(sctx, scerr, maxissues, issues, &wg, resp.LastPage)

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

	dissues := c.processIssues(ctx, cerr, issues, issuescounter, &wgT)

	go func() {
		wgTP.Wait()
		wgT.Wait()
		close(scerr)
		close(issuescounter)
	}()

	return dissues, issuescounter, maxissues
}

// processIssuesPage gets the Issues from GitHub for given page and returns them via
// given channel. IssuesCount contains the amount of issues in the current response
func (c *Connector) processIssuesPage(
	ctx context.Context,
	page int,
	ret chan<- []*github.Issue,
) (
	resp *github.Response,
	issuesCount int,
	err error,
) {

	issues, resp, err := c.client.Issues.ListByRepo(
		ctx,
		c.Owner,
		c.Repo,
		&github.IssueListByRepoOptions{
			State:       "closed",
			ListOptions: github.ListOptions{Page: page, PerPage: IssuesPerPage},
		},
	)
	if err != nil {
		return nil, 0, err
	}

	select {
	case <-ctx.Done():
		return nil, 0, nil
	case ret <- issues:
		return resp, len(issues), nil
	}
}

// processIssuesPages processes GitHub Issues page numbers, given in the cpages channel and returns
// the GH Issues data structures via channel
// possible errors are returned via given cerr channel
func (c *Connector) processIssuesPages(
	ctx context.Context,
	cerr chan<- error,
	cmaxissues chan<- int,
	issues chan<- []*github.Issue,
	wg *sync.WaitGroup,
	lastPage int,
) (cpages chan<- int) {
	ret := make(chan int)

	for i := 0; i < issuesProcessingRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for page := range ret {
				_, n, err := c.processIssuesPage(ctx, page, issues)

				if err != nil {
					helpers.NonBlockingErrSend(ctx, cerr, err)
					return
				}

				if page == lastPage {
					cmaxissues <- n + (lastPage-1)*IssuesPerPage
				}
			}
		}()
	}

	return ret
}

func (c *Connector) processIssues(
	ctx context.Context,
	_ chan<- error,
	cissues <-chan []*github.Issue,
	cissuescounter chan<- bool,
	wg *sync.WaitGroup,
) <-chan data.Issue {

	ret := make(chan data.Issue)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for issues := range cissues {
				for _, issue := range issues {
					cissuescounter <- true
					//ensure we have an issue and not PR
					if issue.PullRequestLinks.GetURL() != "" {
						continue
					}

					var lbs []string
					if issue.Labels != nil && len(issue.Labels) > 0 {
						for _, l := range issue.Labels {
							lbs = append(lbs, *l.Name)
						}
					}

					issue := data.Issue{
						ID:         issue.GetNumber(),
						Name:       issue.GetTitle(),
						ClosedDate: issue.GetClosedAt(),
						URL:        issue.GetHTMLURL(),
						Labels:     lbs,
					}

					select {
					case <-ctx.Done():
						return
					case ret <- issue:
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
