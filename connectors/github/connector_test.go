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
	RetListTags      []*github.RepositoryTag
	RetListTagsErr   error
	RetGetCommits    map[string]*github.RepositoryCommit
	RetGetCommitsErr error
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

func getStringPtr(s string) *string {
	return &s
}

func getTimePtr(t time.Time) *time.Time {
	return &t
}

func Test_connector_GetTags(t *testing.T) {
	type fields struct {
		API   cgithub.API
		Owner string
		Repo  string
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
						&github.RepositoryTag{
							Name: getStringPtr("v0.0.1"),
							Commit: &github.Commit{
								SHA: getStringPtr("7d84cdb2f7c2d4619cda4b8adeb1897097b5c8fc"),
							},
						},
						&github.RepositoryTag{
							Name: getStringPtr("v0.0.2"),
							Commit: &github.Commit{
								SHA: getStringPtr("b3622b516b8ad70ce5dc3fa422fb90c3b58fa9da"),
							},
						},
					},
					RetGetCommits: map[string]*github.RepositoryCommit{
						"7d84cdb2f7c2d4619cda4b8adeb1897097b5c8fc": &github.RepositoryCommit{
							Commit: &github.Commit{
								SHA: getStringPtr("7d84cdb2f7c2d4619cda4b8adeb1897097b5c8fc"),
								Committer: &github.CommitAuthor{
									Date: getTimePtr(time.Unix(2147483647, 0)),
								},
							},
						},
						"b3622b516b8ad70ce5dc3fa422fb90c3b58fa9da": &github.RepositoryCommit{
							Commit: &github.Commit{
								SHA: getStringPtr("b3622b516b8ad70ce5dc3fa422fb90c3b58fa9da"),
								Committer: &github.CommitAuthor{
									Date: getTimePtr(time.Unix(2047483647, 0)),
								},
							},
						},
					},
				},
				Owner: "testowner",
				Repo:  "restrepo",
			},
			want: connectors.Tags{
				{
					Name:   "v0.0.1",
					Commit: "7d84cdb2f7c2d4619cda4b8adeb1897097b5c8fc",
					Date:   time.Unix(2147483647, 0),
				},
				{
					Name:   "v0.0.2",
					Commit: "b3622b516b8ad70ce5dc3fa422fb90c3b58fa9da",
					Date:   time.Unix(2047483647, 0),
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
						&github.RepositoryTag{
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
				API:   tt.fields.API,
				Owner: tt.fields.Owner,
				Repo:  tt.fields.Repo,
			}
			got, err := c.GetTags()
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("Connector.GetTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("connector.GetTags() = %v, want %v", got, tt.want)
			}
		})
	}
}