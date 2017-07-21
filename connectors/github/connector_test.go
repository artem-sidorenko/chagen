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

package github_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/artem-sidorenko/chagen/connectors"
	cgithub "github.com/artem-sidorenko/chagen/connectors/github"
	"github.com/google/go-github/github"
)

type testAPIClient struct {
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
func (t *testAPIClient) ListTags(_, _ string) ([]*github.RepositoryTag, error) {
	if t.RetListTagsErr != nil {
		return nil, t.RetListTagsErr
	}
	return t.RetListTags, nil
}

// GetCommit simlualtes the github.Client.Repositories.GetCommit()
func (t *testAPIClient) GetCommit(_, _, sha string) (*github.RepositoryCommit, error) {
	if t.RetGetCommitsErr != nil {
		return nil, t.RetGetCommitsErr
	}
	return t.RetGetCommits[sha], nil
}

// ListIssues simulates the github.Client.Issues.ListByRepo()
func (t *testAPIClient) ListIssues(_, _ string) ([]*github.Issue, error) {
	if t.RetListIssuesErr != nil {
		return nil, t.RetListIssuesErr
	}
	return t.RetListIssues, nil
}

// ListPRs simulates the github.Client.PullRequests.List()
func (t *testAPIClient) ListPRs(_, _ string) ([]*github.PullRequest, error) {
	if t.RetPRsErr != nil {
		return nil, t.RetPRsErr
	}
	return t.RetPRs, nil
}

// GetReleaseByTag simulates the github.Client.RepositoriesService.GetReleaseByTag
func (t *testAPIClient) GetReleaseByTag(_, _, tag string) (*github.RepositoryRelease, error) {
	if t.RetGetReleaseByTagErr != nil {
		return nil, t.RetGetReleaseByTagErr
	}
	return t.RetGetReleaseByTag[tag], nil
}

func getStringPtr(s string) *string {
	return &s
}

func getIntPtr(i int) *int {
	return &i
}

func getTimePtr(t time.Time) *time.Time {
	return &t
}

func Test_connector_GetTags(t *testing.T) {
	type fields struct {
		API        cgithub.API
		Owner      string
		Repo       string
		ProjectURL string
	}
	tests := []struct {
		name    string
		fields  fields
		want    connectors.Tags
		wantErr error
	}{
		{
			name: "API returns proper data",
			fields: fields{
				API: &testAPIClient{
					RetListTags: []*github.RepositoryTag{
						{
							Name: getStringPtr("v0.0.1"),
							Commit: &github.Commit{
								SHA: getStringPtr("7d84cdb2f7c2d4619cda4b8adeb1897097b5c8fc"),
							},
						},
						{
							Name: getStringPtr("v0.0.2"),
							Commit: &github.Commit{
								SHA: getStringPtr("b3622b516b8ad70ce5dc3fa422fb90c3b58fa9da"),
							},
						},
					},
					RetGetReleaseByTag: map[string]*github.RepositoryRelease{
						"v0.0.1": {
							TagName: getStringPtr("v0.0.1"),
							HTMLURL: getStringPtr("https://example.com/releases/v0.0.1"),
						},
					},
					RetGetCommits: map[string]*github.RepositoryCommit{
						"7d84cdb2f7c2d4619cda4b8adeb1897097b5c8fc": {
							Commit: &github.Commit{
								SHA: getStringPtr("7d84cdb2f7c2d4619cda4b8adeb1897097b5c8fc"),
								Committer: &github.CommitAuthor{
									Date: getTimePtr(time.Unix(2147483647, 0)),
								},
							},
						},
						"b3622b516b8ad70ce5dc3fa422fb90c3b58fa9da": {
							Commit: &github.Commit{
								SHA: getStringPtr("b3622b516b8ad70ce5dc3fa422fb90c3b58fa9da"),
								Committer: &github.CommitAuthor{
									Date: getTimePtr(time.Unix(2047483647, 0)),
								},
							},
						},
					},
				},
				Owner:      "testowner",
				Repo:       "restrepo",
				ProjectURL: "https://example.com/testowner/restrepo",
			},
			want: connectors.Tags{
				{
					Name:   "v0.0.1",
					Commit: "7d84cdb2f7c2d4619cda4b8adeb1897097b5c8fc",
					Date:   time.Unix(2147483647, 0),
					URL:    "https://example.com/releases/v0.0.1",
				},
				{
					Name:   "v0.0.2",
					Commit: "b3622b516b8ad70ce5dc3fa422fb90c3b58fa9da",
					Date:   time.Unix(2047483647, 0),
					URL:    "https://example.com/testowner/restrepo/tree/v0.0.2",
				},
			},
		},
		{
			name: "ListTags call fails",
			fields: fields{
				API: &testAPIClient{
					RetListTagsErr: errors.New("ListTags failed"),
				},
			},
			wantErr: errors.New("ListTags failed"),
		},
		{
			name: "GetCommit call fails",
			fields: fields{
				API: &testAPIClient{
					RetListTags: []*github.RepositoryTag{
						{
							Name: getStringPtr("v0.0.1"),
							Commit: &github.Commit{
								SHA: getStringPtr("7d84cdb2f7c2d4619cda4b8adeb1897097b5c8fc"),
							},
						},
					},
					RetGetCommitsErr: errors.New("GetCommit failed"),
				},
			},
			wantErr: errors.New("GetCommit failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cgithub.Connector{
				API:        tt.fields.API,
				Owner:      tt.fields.Owner,
				Repo:       tt.fields.Repo,
				ProjectURL: tt.fields.ProjectURL,
			}
			got, err := c.GetTags()
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("Connector.GetTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connector.GetTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConnector_GetIssues(t *testing.T) {
	type fields struct {
		API   cgithub.API
		Owner string
		Repo  string
	}
	tests := []struct {
		name    string
		fields  fields
		want    connectors.Issues
		wantErr error
	}{
		{
			name: "API returns proper data",
			fields: fields{
				API: &testAPIClient{
					RetListIssues: []*github.Issue{
						{
							Number:           getIntPtr(1234),
							Title:            getStringPtr("Test issue title"),
							PullRequestLinks: &github.PullRequestLinks{},
							ClosedAt:         getTimePtr(time.Unix(1047483647, 0)),
							HTMLURL:          getStringPtr("http://example.com/issues/1234"),
						},
						{
							Number: getIntPtr(4321),
							Title:  getStringPtr("Test PR title"),
							PullRequestLinks: &github.PullRequestLinks{
								URL: getStringPtr("https://example.com/prs/4321"),
							},
						},
					},
				},
			},
			want: connectors.Issues{
				connectors.Issue{
					ID:         1234,
					Name:       "Test issue title",
					ClosedDate: time.Unix(1047483647, 0),
					URL:        "http://example.com/issues/1234",
				},
			},
		},
		{
			name: "ListIssues call fails",
			fields: fields{
				API: &testAPIClient{
					RetListIssuesErr: errors.New("ListIssues failed"),
				},
			},
			wantErr: errors.New("ListIssues failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cgithub.Connector{
				API:   tt.fields.API,
				Owner: tt.fields.Owner,
				Repo:  tt.fields.Repo,
			}
			got, err := c.GetIssues()
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("Connector.GetIssues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connector.GetIssues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConnector_GetMRs(t *testing.T) {
	type fields struct {
		API   cgithub.API
		Owner string
		Repo  string
	}
	tests := []struct {
		name    string
		fields  fields
		want    connectors.MRs
		wantErr error
	}{
		{
			name: "API returns proper data",
			fields: fields{
				API: &testAPIClient{
					RetPRs: []*github.PullRequest{
						{
							Number:  getIntPtr(1234),
							Title:   getStringPtr("Test PR title"),
							HTMLURL: getStringPtr("https://example.com/pulls/1234"),
							User: &github.User{
								Login:   getStringPtr("test-user"),
								HTMLURL: getStringPtr("https://example.com/users/test-user"),
							},
							MergedAt: getTimePtr(time.Unix(1747483647, 0)),
						},
						{
							Number: getIntPtr(1233),
							Title:  getStringPtr("Second closed PR title"),
						},
					},
				},
			},
			want: connectors.MRs{
				connectors.MR{
					ID:         1234,
					Name:       "Test PR title",
					URL:        "https://example.com/pulls/1234",
					Author:     "test-user",
					AuthorURL:  "https://example.com/users/test-user",
					MergedDate: time.Unix(1747483647, 0),
				},
			},
		},
		{
			name: "ListPRs call fails",
			fields: fields{
				API: &testAPIClient{
					RetPRsErr: errors.New("ListPRs failed"),
				},
			},
			wantErr: errors.New("ListPRs failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cgithub.Connector{
				API:   tt.fields.API,
				Owner: tt.fields.Owner,
				Repo:  tt.fields.Repo,
			}
			got, err := c.GetMRs()
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("Connector.GetMRs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connector.GetMRs() = %v, want %v", got, tt.want)
			}
		})
	}
}
