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

package testapiclient

import (
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
