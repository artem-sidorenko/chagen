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

// TagsPerPage defined how many tags are fetched per page
var TagsPerPage = 30 // nolint: gochecknoglobals

const (
	tagProcessingRoutines = 10
)

// Tags returns the git tags via channels.
// Returns possible errors via given cerr channel
// ctags returns tags
// cmaxtags returns the max available amount of tags
func (c *Connector) Tags(
	ctx context.Context,
	cerr chan<- error,
) (
	ctags <-chan data.Tag,
	cmaxtags <-chan int,
) {
	tags := make(chan []*github.RepositoryTag)
	maxtags := make(chan int)

	// establshing the local error handling with local context
	// if we get any local errors, cancel all local running routines,
	// close channels and pass the error to the caller
	sctx, cancel := context.WithCancel(ctx)
	scerr := make(chan error)

	go func() {
		select {
		case <-ctx.Done():
			return
		case err := <-scerr:
			cancel()
			cerr <- err
		}
	}()

	// process the list of tags in the background
	// this is a bit special and complex:
	// - we have to request the first page of information in order to know the
	//   amount of data we would get
	// - we do not throw this page away, but we completely use the data delivered
	// - to be DRY we have to use processTagPage, so we have to use the
	//   tags channel in processTagPage and processTagPages,
	//   so we can't close the tags channel in processTagPages
	// - we share the wg sync group in order to know when all go routines
	//   are finished and we can close the channels
	go func() {
		var wg sync.WaitGroup

		// we use this func to close all related channels and invoke it in two places:
		// - we might have to close them early, where no sub goroutines are invoked
		//   (see the first data fetch and err handling)
		// - we close them after all goroutines are finished
		closeCh := func() {
			close(tags)
			close(maxtags)
		}

		// get the first page in order to know the amount of data
		resp, n, err := c.processTagPage(sctx, 1, tags)
		if err != nil {
			scerr <- err
			closeCh()
			return
		}

		if resp.LastPage == 0 { // if we have only one page, we already know the amount of data
			maxtags <- n
		} else {
			// spawn goroutines for page processing
			cpages := c.processTagPages(sctx, scerr, maxtags, tags, &wg, resp.LastPage)
			// spawn for each page an own go routine
			// we start from the last page as we want to get the max amount of data fast
			for i := resp.LastPage; i >= 2; i-- {
				cpages <- i
			}
			close(cpages)
		}

		go func() { // all processing routines are done -> close the channels
			wg.Wait()
			closeCh()
		}()
	}()

	dtags := c.processTags(sctx, scerr, tags)

	return dtags, maxtags
}

// processTagPage gets the tags from GitHub for given page and returns them via
// given channel. tagCount contains the amount of tags in the current response
func (c *Connector) processTagPage(
	ctx context.Context,
	page int,
	ret chan<- []*github.RepositoryTag,
) (
	resp *github.Response,
	tagsCount int,
	err error,
) {
	tags, resp, err := c.client.Repositories.ListTags(
		ctx,
		c.Owner,
		c.Repo,
		&github.ListOptions{Page: page, PerPage: TagsPerPage},
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

// processTagPages processes  GitHub tag page numbers, given in the cpages channel and returns
// the GH RepositoryTag data structures via channel
// possible errors are returned via given cerr channel
func (c *Connector) processTagPages(
	ctx context.Context,
	cerr chan<- error,
	cmaxtags chan<- int,
	tags chan<- []*github.RepositoryTag,
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
					cerr <- err
					return
				}

				if page == lastPage { // this is the last page, lets calculate amount of data
					cmaxtags <- n + (lastPage-1)*TagsPerPage
				}
			}
		}()
	}

	return ret
}

// processTags processes given GitHub tags in the ctags channel and returns
// the tags in our data structure via channel
// possible errors are returned via given cerr channel
func (c *Connector) processTags(
	ctx context.Context,
	cerr chan<- error,
	ctags <-chan []*github.RepositoryTag,
) <-chan data.Tag {

	ret := make(chan data.Tag)
	var wg sync.WaitGroup

	for i := 0; i < tagProcessingRoutines; i++ {
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
