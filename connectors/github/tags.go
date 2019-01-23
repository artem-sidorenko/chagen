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

// listTags gets the tags from GitHub client and returns them via channel
// possible errors are returned via given cerr channel
func (c *Connector) listTags(ctx context.Context, cerr chan error) <-chan []*github.RepositoryTag {
	ret := make(chan []*github.RepositoryTag)

	go func() {
		defer close(ret)

		opt := &github.ListOptions{}
		for {
			tags, resp, err := c.client.Repositories.ListTags(ctx, c.Owner, c.Repo, opt)
			if err != nil {
				cerr <- err
				return
			}

			select {
			case <-ctx.Done():
				return
			case ret <- tags:
			}

			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
		}
	}()

	return ret
}

// processTags processes given GitHub tags in the ctags channel and returns
// the tags in our data structure via channel
// possible errors are returned via given cerr channel
func (c *Connector) processTags(
	ctx context.Context,
	cerr chan error,
	ctags <-chan []*github.RepositoryTag,
) <-chan data.Tag {

	ret := make(chan data.Tag)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for tags := range ctags {
				for _, tag := range tags {
					tagName := tag.GetName()

					commit, _, err := c.client.Repositories.GetCommit(ctx,
						c.Owner, c.Repo, tag.Commit.GetSHA())
					if err != nil {
						cerr <- err
						return
					}

					tagURL, err := c.getTagURL(tagName, false)
					if err != nil {
						cerr <- err
						return
					}

					tag := data.Tag{
						Name:   tagName,
						Commit: commit.Commit.GetSHA(),
						Date:   commit.Commit.Committer.GetDate(),
						URL:    tagURL,
					}

					select {
					case <-ctx.Done():
						return
					case ret <- tag:
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
