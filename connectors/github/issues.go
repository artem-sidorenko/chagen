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

	"github.com/artem-sidorenko/chagen/data"
	"github.com/google/go-github/github"
)

// Issues returns the issues via channels.
// Returns possible errors via given cerr channel
// cissues returns issues
// cmaxissues returns the max available amount of issues
func (c *Connector) Issues(
	ctx context.Context,
	cerr chan<- error,
) (
	ctags <-chan data.Issue,
	cmaxissues <-chan int,
) {
	issues := c.listIssues(ctx, cerr)
	dissues := c.processIssues(ctx, cerr, issues)

	return dissues, nil
}

func (c *Connector) listIssues(
	ctx context.Context,
	cerr chan<- error,
) <-chan []*github.Issue {
	ret := make(chan []*github.Issue)

	go func() {
		defer close(ret)

		opt := &github.IssueListByRepoOptions{State: "closed"}

		for {
			issues, resp, err := c.client.Issues.ListByRepo(ctx, c.Owner, c.Repo, opt)
			if err != nil {
				cerr <- err
				return
			}

			select {
			case <-ctx.Done():
				return
			case ret <- issues:
			}

			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
		}
	}()

	return ret
}

func (c *Connector) processIssues(
	ctx context.Context,
	_ chan<- error,
	cissues <-chan []*github.Issue,
) <-chan data.Issue {

	ret := make(chan data.Issue)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for issues := range cissues {
				for _, issue := range issues {
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
