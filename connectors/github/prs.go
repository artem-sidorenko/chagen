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

	"github.com/artem-sidorenko/chagen/data"
	"github.com/google/go-github/github"
)

// MRs returns the PRs via channels.
// Returns possible errors via given cerr channel
// cmrs returns MRs
// cmaxmrs returns the max available amount of MRs
func (c *Connector) MRs(
	ctx context.Context,
	cerr chan<- error,
	cmaxmrs chan<- int,
) (
	cmrs <-chan data.MR,
) {
	mrs := c.listPRs(ctx, cerr)
	dmrs := c.processPRs(ctx, cerr, mrs)

	return dmrs
}

func (c *Connector) listPRs(
	ctx context.Context,
	cerr chan<- error,
) <-chan []*github.PullRequest {
	ret := make(chan []*github.PullRequest)

	go func() {
		defer close(ret)

		opt := &github.PullRequestListOptions{State: "closed"}

		for {
			prs, resp, err := c.client.PullRequests.List(ctx, c.Owner, c.Repo, opt)
			if err != nil {
				cerr <- err
				return
			}

			select {
			case <-ctx.Done():
				return
			case ret <- prs:
			}

			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
		}
	}()

	return ret
}

func (c *Connector) processPRs(
	ctx context.Context,
	_ chan<- error,
	cprs <-chan []*github.PullRequest,
) <-chan data.MR {

	ret := make(chan data.MR)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for prs := range cprs {
				for _, pr := range prs {
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
