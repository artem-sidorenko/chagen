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

// Package testapiclient implements the used interfaces of github API client library
// and simulates the API answers for our tests
package testapiclient

import (
	"errors"
	"time"

	"github.com/google/go-github/github"
)

// TestAPIClient simulates the responses of github API client library
type TestAPIClient struct {
	RetListTags           []*github.RepositoryTag
	RetListTagsErr        error
	RetGetCommits         map[string]*github.RepositoryCommit
	RetGetCommitsErr      error
	RetListIssues         []*github.Issue
	RetListIssuesErr      error
	RetPRs                []*github.PullRequest
	RetPRsErr             error
	RetGetReleaseByTag    map[string]*github.RepositoryRelease
	RetGetReleaseByTagErr error
}

// ListTags simulates the github.Client.Repositories.ListTags()
func (t *TestAPIClient) ListTags(_, _ string) ([]*github.RepositoryTag, error) {
	if t.RetListTagsErr != nil {
		return nil, t.RetListTagsErr
	}
	return t.RetListTags, nil
}

// GetCommit simlualtes the github.Client.Repositories.GetCommit()
func (t *TestAPIClient) GetCommit(_, _, sha string) (*github.RepositoryCommit, error) {
	if t.RetGetCommitsErr != nil {
		return nil, t.RetGetCommitsErr
	}
	return t.RetGetCommits[sha], nil
}

// ListIssues simulates the github.Client.Issues.ListByRepo()
func (t *TestAPIClient) ListIssues(_, _ string) ([]*github.Issue, error) {
	if t.RetListIssuesErr != nil {
		return nil, t.RetListIssuesErr
	}
	return t.RetListIssues, nil
}

// ListPRs simulates the github.Client.PullRequests.List()
func (t *TestAPIClient) ListPRs(_, _ string) ([]*github.PullRequest, error) {
	if t.RetPRsErr != nil {
		return nil, t.RetPRsErr
	}
	return t.RetPRs, nil
}

// GetReleaseByTag simulates the github.Client.RepositoriesService.GetReleaseByTag
func (t *TestAPIClient) GetReleaseByTag(_, _, tag string) (*github.RepositoryRelease, error) {
	if t.RetGetReleaseByTagErr != nil {
		return nil, t.RetGetReleaseByTagErr
	}
	return t.RetGetReleaseByTag[tag], nil
}

// GetStringPtr returns a pointer for a given string
func GetStringPtr(s string) *string {
	return &s
}

// GetIntPtr returns a pointer for a given int
func GetIntPtr(i int) *int {
	return &i
}

// GetTimePtr returns a pointer for a given Time
func GetTimePtr(t time.Time) *time.Time {
	return &t
}

// Options represents the possible options for API simulation
// if a field is set to true - return valid data, if it set to false - error
type Options struct {
	ListTags   bool
	GetCommit  bool
	ListIssues bool
	ListPRs    bool
}

// New returns a TestAPIClient with according data or errors
func New(o Options) *TestAPIClient {
	t := &TestAPIClient{}

	if o.ListTags {
		t.RetListTags = []*github.RepositoryTag{
			{
				Name: GetStringPtr("v0.0.1"),
				Commit: &github.Commit{
					SHA: GetStringPtr("7d84cdb2f7c2d4619cda4b8adeb1897097b5c8fc"),
				},
			},
			{
				Name: GetStringPtr("v0.0.2"),
				Commit: &github.Commit{
					SHA: GetStringPtr("b3622b516b8ad70ce5dc3fa422fb90c3b58fa9da"),
				},
			},
		}
		t.RetGetReleaseByTag = map[string]*github.RepositoryRelease{
			"v0.0.1": {
				TagName: GetStringPtr("v0.0.1"),
				HTMLURL: GetStringPtr("https://example.com/releases/v0.0.1"),
			},
		}
	} else {
		t.RetListTagsErr = errors.New("ListTags failed")
	}

	if o.GetCommit {
		t.RetGetCommits = map[string]*github.RepositoryCommit{
			"7d84cdb2f7c2d4619cda4b8adeb1897097b5c8fc": {
				Commit: &github.Commit{
					SHA: GetStringPtr("7d84cdb2f7c2d4619cda4b8adeb1897097b5c8fc"),
					Committer: &github.CommitAuthor{
						Date: GetTimePtr(time.Unix(2147483647, 0)),
					},
				},
			},
			"b3622b516b8ad70ce5dc3fa422fb90c3b58fa9da": {
				Commit: &github.Commit{
					SHA: GetStringPtr("b3622b516b8ad70ce5dc3fa422fb90c3b58fa9da"),
					Committer: &github.CommitAuthor{
						Date: GetTimePtr(time.Unix(2047483647, 0)),
					},
				},
			},
		}
	} else {
		t.RetGetCommitsErr = errors.New("GetCommit failed")
	}

	if o.ListIssues {
		t.RetListIssues = []*github.Issue{
			{
				Number:           GetIntPtr(1234),
				Title:            GetStringPtr("Test issue title"),
				PullRequestLinks: &github.PullRequestLinks{},
				ClosedAt:         GetTimePtr(time.Unix(1047483647, 0)),
				HTMLURL:          GetStringPtr("http://example.com/issues/1234"),
			},
			{
				Number: GetIntPtr(4321),
				Title:  GetStringPtr("Test PR title"),
				PullRequestLinks: &github.PullRequestLinks{
					URL: GetStringPtr("https://example.com/prs/4321"),
				},
			},
		}
	} else {
		t.RetListIssuesErr = errors.New("ListIssues failed")
	}

	if o.ListPRs {
		t.RetPRs = []*github.PullRequest{
			{
				Number:  GetIntPtr(1234),
				Title:   GetStringPtr("Test PR title"),
				HTMLURL: GetStringPtr("https://example.com/pulls/1234"),
				User: &github.User{
					Login:   GetStringPtr("test-user"),
					HTMLURL: GetStringPtr("https://example.com/users/test-user"),
				},
				MergedAt: GetTimePtr(time.Unix(1747483647, 0)),
			},
			{
				Number: GetIntPtr(1233),
				Title:  GetStringPtr("Second closed PR title"),
			},
		}
	} else {
		t.RetPRsErr = errors.New("ListPRs failed")
	}

	return t
}
