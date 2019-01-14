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

	cgithub "github.com/artem-sidorenko/chagen/connectors/github"
	"github.com/artem-sidorenko/chagen/connectors/github/internal/testapiclient"
	"github.com/artem-sidorenko/chagen/data"
	"github.com/google/go-github/github"
)

// NewConnector returns a new Connector, initialited with test data
// and provided test options
func NewConnector(o testapiclient.Options) *cgithub.Connector {
	return &cgithub.Connector{
		API:        testapiclient.New(o),
		Owner:      "testowner",
		Repo:       "restrepo",
		ProjectURL: "https://example.com/testowner/restrepo",
	}
}

func Test_connector_GetTags(t *testing.T) {
	type fields struct {
		TestAPIopts testapiclient.Options
	}
	tests := []struct {
		name    string
		fields  fields
		want    data.Tags
		wantErr error
	}{
		{
			name: "API returns proper data",
			fields: fields{
				TestAPIopts: testapiclient.Options{
					ListTags:  true,
					GetCommit: true,
				},
			},
			want: data.Tags{
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
				TestAPIopts: testapiclient.Options{ListTags: false},
			},
			wantErr: errors.New("ListTags failed"),
		},
		{
			name: "GetCommit call fails",
			fields: fields{
				TestAPIopts: testapiclient.Options{ListTags: true, GetCommit: false},
			},
			wantErr: errors.New("GetCommit failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConnector(tt.fields.TestAPIopts)
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
		TestAPIopts testapiclient.Options
	}
	tests := []struct {
		name    string
		fields  fields
		want    data.Issues
		wantErr error
	}{
		{
			name: "API returns proper data",
			fields: fields{
				TestAPIopts: testapiclient.Options{ListIssues: true},
			},
			want: data.Issues{
				data.Issue{
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
				TestAPIopts: testapiclient.Options{ListIssues: false},
			},
			wantErr: errors.New("ListIssues failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConnector(tt.fields.TestAPIopts)
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
		TestAPIopts testapiclient.Options
	}
	tests := []struct {
		name    string
		fields  fields
		want    data.MRs
		wantErr error
	}{
		{
			name: "API returns proper data",
			fields: fields{
				TestAPIopts: testapiclient.Options{ListPRs: true},
			},
			want: data.MRs{
				data.MR{
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
				TestAPIopts: testapiclient.Options{ListPRs: false},
			},
			wantErr: errors.New("ListPRs failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConnector(tt.fields.TestAPIopts)
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

func TestConnector_GetNewTagURL(t *testing.T) {
	type fields struct {
		NewTagUseReleaseURL bool
	}
	type args struct {
		TagName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Release is present, NewTagUseReleaseURL is disabled",
			fields: fields{
				NewTagUseReleaseURL: false,
			},
			args: args{
				TagName: "v0.0.1",
			},
			want: "https://example.com/testowner/restrepo/releases/v0.0.1",
		},
		{
			name: "Release is present, NewTagUseReleaseURL is enabled",
			fields: fields{
				NewTagUseReleaseURL: true,
			},
			args: args{
				TagName: "v0.0.1",
			},
			want: "https://example.com/testowner/restrepo/releases/v0.0.1",
		},
		{
			name: "Release is not present, NewTagUseReleaseURL is disabled",
			fields: fields{
				NewTagUseReleaseURL: false,
			},
			args: args{
				TagName: "v0.0.3",
			},
			want: "https://example.com/testowner/restrepo/tree/v0.0.3",
		},
		{
			name: "Release is not present, alwaysUseReleaseURL is enabled",
			fields: fields{
				NewTagUseReleaseURL: true,
			},
			args: args{
				TagName: "v0.0.3",
			},
			want: "https://example.com/testowner/restrepo/releases/v0.0.3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cgithub.Connector{
				API: &testapiclient.TestAPIClient{
					RetGetReleaseByTag: map[string]*github.RepositoryRelease{
						"v0.0.1": {
							TagName: testapiclient.GetStringPtr("v0.0.1"),
							HTMLURL: testapiclient.GetStringPtr(
								"https://example.com/testowner/restrepo/releases/v0.0.1",
							),
						},
					},
				},
				ProjectURL:          "https://example.com/testowner/restrepo",
				NewTagUseReleaseURL: tt.fields.NewTagUseReleaseURL,
			}
			got, err := c.GetNewTagURL(tt.args.TagName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connector.GetNewTagURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Connector.GetNewTagURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
