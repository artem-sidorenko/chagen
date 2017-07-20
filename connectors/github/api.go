/*
   Copyright 2017 Artem Sidorenko <artem@posteo.de>

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

	"fmt"

	"github.com/google/go-github/github"
)

// API builds a wrapper around GitHub API from
// github.com/google/go-github/github
// We need this in order to allow unit testing without to fire
// requests to GitHub during the tests
type API interface {
	ListTags(string, string) ([]*github.RepositoryTag, error)
	GetCommit(string, string, string) (*github.RepositoryCommit, error)
	ListIssues(string, string) ([]*github.Issue, error)
}

// APIClient implements the API interface
// and connects the needed functions to the github API library
type APIClient struct {
	context context.Context
	client  *github.Client
}

// ListTags implements the github.Client.Repositories.ListTags()
func (a *APIClient) ListTags(owner, repo string) (tags []*github.RepositoryTag, err error) {
	tags, _, err = a.client.Repositories.ListTags(a.context, owner, repo, nil)
	if err != nil {
		return nil, a.formatErrorCode("ListTags", err)
	}

	return
}

// GetCommit implements the github.Client.Repositories.GetCommit()
func (a *APIClient) GetCommit(owner, repo, sha string) (
	commit *github.RepositoryCommit, err error) {

	commit, _, err = a.client.Repositories.GetCommit(a.context, owner, repo, sha)
	if err != nil {
		return nil, a.formatErrorCode("GetCommit", err)
	}

	return
}

// ListIssues implements the github.Client.Issues.ListByRepo()
// and returns all closed issues and PRs
func (a *APIClient) ListIssues(owner, repo string) ([]*github.Issue, error) {
	opt := &github.IssueListByRepoOptions{
		State: "closed",
	}

	var ret []*github.Issue
	for {
		issues, resp, err := a.client.Issues.ListByRepo(a.context, owner, repo, opt)
		if err != nil {
			return nil, a.formatErrorCode("ListIssues", err)
		}

		ret = append(ret, issues...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return ret, nil
}

// formatErrorCode formats the error message for this connector
func (a *APIClient) formatErrorCode(query string, err error) error {
	return fmt.Errorf("GitHub query '%s' failed: %s", query, err)
}

// NewAPIClient returns the initialized and ready to use APIClient
func NewAPIClient() *APIClient {
	return &APIClient{
		context: context.Background(),
		client:  github.NewClient(nil),
	}
}
