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

	"github.com/artem-sidorenko/chagen/data"
	gitlab "github.com/xanzy/go-gitlab"
)

// TagsPerPage defined how many tags are fetched per page
var TagsPerPage = 30 // nolint: gochecknoglobals

const (
	tagProcessingRoutines = 10
)

// Tags returns the git tags via channels.
// Returns possible errors via given cerr channel
// ctags returns tags
// ctagscounter returns the channel, which ticks when a tag is proceeded
// cmaxtags returns the max available amount of tags
func (c *Connector) Tags(
	ctx context.Context,
	cerr chan<- error,
) (
	ctags <-chan data.Tag,
	ctagscounter <-chan bool,
	cmaxtags <-chan int,
) {
	// for detailed comments, please see the github/tags.go
	tags := make(chan []*gitlab.Tag)
	maxtags := make(chan int)
	tagscounter := make(chan bool, 100)

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
			close(tags)
			close(maxtags)
		}

		resp, n, err := c.processTagPage(sctx, 1, tags)
		if err != nil {
			nonBlockingErrSend(sctx, scerr, err)
			closeCh()
			return
		}

		if resp.TotalPages == 1 {
			maxtags <- n
		} else {
			cpages := c.processTagPages(sctx, scerr, maxtags, tags, &wg, resp.TotalPages)

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

	dtags := c.processTags(sctx, scerr, tags, tagscounter, &wgT)

	go func() {
		wgTP.Wait()
		wgT.Wait()
		close(scerr)
	}()

	return dtags, tagscounter, maxtags
}

// processTagPage gets the tags from GitLab for given page and returns them via
// given channel. tagCount contains the amount of tags in the current response
func (c *Connector) processTagPage(
	ctx context.Context,
	page int,
	ret chan<- []*gitlab.Tag,
) (
	resp *gitlab.Response,
	tagsCount int,
	err error,
) {
	tags, resp, err := c.client.Tags.ListTags(
		c.Owner+"/"+c.Repo,
		&gitlab.ListTagsOptions{
			ListOptions: gitlab.ListOptions{Page: page, PerPage: TagsPerPage}},
	)

	if err != nil {
		return nil, 0, err
	}

	select {
	case <-ctx.Done():
		return nil, 0, nil
	case ret <- tags:
		return resp, len(tags), nil
	}
}

// processTagPages processes GitLab tag page numbers, given in the cpages channel and returns
// the gitlab.Tag data structures via channel
// possible errors are returned via given cerr channel
func (c *Connector) processTagPages(
	ctx context.Context,
	cerr chan<- error,
	cmaxtags chan<- int,
	tags chan<- []*gitlab.Tag,
	wg *sync.WaitGroup,
	lastPage int,
) (cpages chan<- int) {
	ret := make(chan int)

	for i := 0; i < tagProcessingRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for page := range ret {
				_, n, err := c.processTagPage(ctx, page, tags)
				if err != nil {
					nonBlockingErrSend(ctx, cerr, err)
					return
				}

				if page == lastPage {
					cmaxtags <- n + (lastPage-1)*TagsPerPage
				}
			}
		}()
	}

	return ret
}

// processTags processes given GitLab tags in the ctags channel and returns
// the tags in our data structure via channel
// possible errors are returned via given cerr channel
func (c *Connector) processTags(
	ctx context.Context,
	cerr chan<- error,
	ctags <-chan []*gitlab.Tag,
	ctagscounter chan<- bool,
	wg *sync.WaitGroup,
) <-chan data.Tag {

	ret := make(chan data.Tag)

	for i := 0; i < tagProcessingRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for tags := range ctags {
				for _, tag := range tags {
					ctagscounter <- true
					tagName := tag.Name

					commit := tag.Commit

					tagURL, err := c.getTagURL(tagName)
					if err != nil {
						nonBlockingErrSend(ctx, cerr, err)
						return
					}

					tag := data.Tag{
						Name:   tagName,
						Commit: commit.ID,
						Date:   *commit.AuthoredDate,
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
