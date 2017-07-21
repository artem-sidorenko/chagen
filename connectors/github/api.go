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

	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// API builds a wrapper around GitHub API from
// github.com/google/go-github/github
// We need this in order to allow unit testing without to fire
// requests to GitHub during the tests
type API interface {
	ListTags(string, string) ([]*github.RepositoryTag, error)
	GetCommit(string, string, string) (*github.RepositoryCommit, error)
	ListIssues(string, string) ([]*github.Issue, error)
	ListPRs(string, string) ([]*github.PullRequest, error)
	GetReleaseByTag(string, string, string) (*github.RepositoryRelease, error)
}

// APIClient implements the API interface
// and connects the needed functions to the github API library
type APIClient struct {
	context context.Context
	client  *github.Client
}

// ListTags implements the github.Client.Repositories.ListTags()
func (a *APIClient) ListTags(owner, repo string) ([]*github.RepositoryTag, error) {
	opt := &github.ListOptions{}

	var ret []*github.RepositoryTag
	for {
		tags, resp, err := a.client.Repositories.ListTags(a.context, owner, repo, nil)
		if err != nil {
			return nil, a.formatErrorCode("ListTags", err)
		}

		ret = append(ret, tags...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return ret, nil
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

// ListPRs implements the github.Client.PullRequests.List()
// and returns all closed PRs
func (a *APIClient) ListPRs(owner, repo string) ([]*github.PullRequest, error) {
	opt := &github.PullRequestListOptions{
		State: "closed",
	}

	var ret []*github.PullRequest
	for {
		prs, resp, err := a.client.PullRequests.List(a.context, owner, repo, opt)
		if err != nil {
			return nil, a.formatErrorCode("ListPRs", err)
		}

		ret = append(ret, prs...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return ret, nil
}

// GetReleaseByTag implements the github.Client.RepositoriesService.GetReleaseByTag
// and returns a release if possible. Returns nil if no release is found
func (a *APIClient) GetReleaseByTag(owner, repo, tag string) (*github.RepositoryRelease, error) {
	release, resp, err := a.client.Repositories.GetReleaseByTag(a.context, owner, repo, tag)
	if err != nil {
		// no release was found for this tag, this is no error for us
		if resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, a.formatErrorCode("GetReleaseByTag", err)
	}

	return release, nil
}

// formatErrorCode formats the error message for this connector
func (a *APIClient) formatErrorCode(query string, err error) error {
	return fmt.Errorf("GitHub query '%s' failed: %s", query, err)
}

// NewAPIClient returns the initialized and ready to use APIClient
// Uses AccessToken for authentication if not empty
func NewAPIClient(AccessToken string) *APIClient {
	ctx := context.Background()

	var tc *http.Client

	if AccessToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: AccessToken},
		)
		tc = oauth2.NewClient(ctx, ts)
	}

	return &APIClient{
		context: ctx,
		client:  github.NewClient(tc),
	}
}
