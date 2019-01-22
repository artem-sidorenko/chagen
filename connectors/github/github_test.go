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
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/artem-sidorenko/chagen/connectors"
	"github.com/artem-sidorenko/chagen/connectors/github"
	"github.com/artem-sidorenko/chagen/connectors/github/internal/testclient"
	"github.com/artem-sidorenko/chagen/data"
	tcli "github.com/artem-sidorenko/chagen/internal/testing/cli"

	"github.com/urfave/cli"
)

func setupTestConnector(
	returnValue testclient.ReturnValueStr, newTagUseReleaseURL bool,
) connectors.Connector {

	github.NewGitHubClientFunc = testclient.New
	cliFlags := map[string]string{
		"github-owner": "testowner",
		"github-repo":  "testrepo",
	}
	if newTagUseReleaseURL {
		cliFlags["github-release-url"] = "true"
	}

	ctx := tcli.TestContext(github.CLIFlags(), cliFlags)

	// initialize error values
	testclient.ReturnValue = returnValue

	c, err := github.New(ctx)
	if err != nil {
		panic(fmt.Sprintf("got error from test consturctor: %v", err))
	}

	return c
}

func TestConnector_RepositoryExists(t *testing.T) {
	tests := []struct {
		name        string
		returnValue testclient.ReturnValueStr
		want        bool
		wantErr     error
	}{
		{
			name: "API returns 200 for Ok",
			returnValue: testclient.ReturnValueStr{
				RetRepoServiceGetRespCode: 200,
			},
			want: true,
		},
		{
			name: "API returns 404",
			returnValue: testclient.ReturnValueStr{
				RetRepoServiceGetRespCode: 404,
			},
			want: false,
		},
		{
			name: "API returns unhandled error code 500",
			returnValue: testclient.ReturnValueStr{
				RetRepoServiceGetRespCode: 500,
			},
			wantErr: errors.New("GitHub query 'RepositoryExists' failed: unhandled HTTP response code 500"),
		},
		{
			name: "Get returns an error",
			returnValue: testclient.ReturnValueStr{
				RetRepoServiceGetErr: true,
			},
			wantErr: errors.New("GitHub query 'RepositoryExists' failed: can't fetch the repo data"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestConnector(tt.returnValue, false)

			got, err := c.RepositoryExists()

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Connector.RepositoryExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Connector.RepositoryExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTags(t *testing.T) {
	tests := []struct {
		name        string
		returnValue testclient.ReturnValueStr
		want        data.Tags
		wantErr     error
	}{
		{
			name: "API returns proper data",
			want: data.Tags{
				{
					Name:   "v0.0.1",
					Commit: "7d84cdb2f7c2d4619cda4b8adeb1897097b5c8fc",
					Date:   time.Unix(2147483647, 0),
					URL:    "https://github.com/testowner/testrepo/releases/v0.0.1",
				},
				{
					Name:   "v0.0.2",
					Commit: "b3622b516b8ad70ce5dc3fa422fb90c3b58fa9da",
					Date:   time.Unix(2047483647, 0),
					URL:    "https://github.com/testowner/testrepo/tree/v0.0.2",
				},
			},
		},
		{
			name: "ListTags call fails",
			returnValue: testclient.ReturnValueStr{
				RetRepoServiceListTagsErr: true,
			},
			wantErr: errors.New("GitHub query 'GetTags' failed: can't fetch the tags"),
		},
		{
			name: "GetCommit call fails",
			returnValue: testclient.ReturnValueStr{
				RetRepoServiceGetCommitsErr: true,
			},
			wantErr: errors.New("GitHub query 'GetTags' failed: can't fetch the commit"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestConnector(tt.returnValue, false)

			got, err := c.GetTags()

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Connector.GetTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connector.GetTags() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestGetIssues(t *testing.T) {
	tests := []struct {
		name        string
		returnValue testclient.ReturnValueStr
		want        data.Issues
		wantErr     error
	}{
		{
			name: "API returns proper data",
			want: data.Issues{
				data.Issue{
					ID:         1234,
					Name:       "Test issue title",
					ClosedDate: time.Unix(1047483647, 0),
					URL:        "http://example.com/issues/1234",
					Labels:     []string{"enhancement"},
				},
			},
		},
		{
			name: "ListIssues call fails",
			returnValue: testclient.ReturnValueStr{
				RetIssueServiceListByRepoErr: true,
			},
			wantErr: errors.New("GitHub query 'GetIssues' failed: can't fetch the issues"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestConnector(tt.returnValue, false)

			got, err := c.GetIssues()
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Connector.GetIssues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connector.GetIssues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMRs(t *testing.T) {
	tests := []struct {
		name        string
		returnValue testclient.ReturnValueStr
		want        data.MRs
		wantErr     error
	}{
		{
			name: "API returns proper data",
			want: data.MRs{
				data.MR{
					ID:         1234,
					Name:       "Test PR title",
					URL:        "https://example.com/pulls/1234",
					Author:     "test-user",
					AuthorURL:  "https://example.com/users/test-user",
					MergedDate: time.Unix(1747483647, 0),
					Labels:     []string{"bugfix"},
				},
			},
		},
		{
			name: "ListPRs call fails",
			returnValue: testclient.ReturnValueStr{
				RetPullRequestsListErr: true,
			},
			wantErr: errors.New("GitHub query 'GetMRs' failed: can't fetch the PRs"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestConnector(tt.returnValue, false)

			got, err := c.GetMRs()
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Connector.GetMRs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connector.GetMRs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNewTagURL(t *testing.T) {
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
			want: "https://github.com/testowner/testrepo/releases/v0.0.1",
		},
		{
			name: "Release is present, NewTagUseReleaseURL is enabled",
			fields: fields{
				NewTagUseReleaseURL: true,
			},
			args: args{
				TagName: "v0.0.1",
			},
			want: "https://github.com/testowner/testrepo/releases/v0.0.1",
		},
		{
			name: "Release is not present, NewTagUseReleaseURL is disabled",
			fields: fields{
				NewTagUseReleaseURL: false,
			},
			args: args{
				TagName: "v0.0.3",
			},
			want: "https://github.com/testowner/testrepo/tree/v0.0.3",
		},
		{
			name: "Release is not present, alwaysUseReleaseURL is enabled",
			fields: fields{
				NewTagUseReleaseURL: true,
			},
			args: args{
				TagName: "v0.0.3",
			},
			want: "https://github.com/testowner/testrepo/releases/v0.0.3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestConnector(testclient.ReturnValueStr{}, tt.fields.NewTagUseReleaseURL)

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

func setupCLIContext(githubOwner, githubRepo bool) *cli.Context {
	flags := map[string]string{}

	if githubOwner {
		flags["github-owner"] = "testowner"
	}

	if githubRepo {
		flags["github-repo"] = "testrepo"
	}

	return tcli.TestContext(github.CLIFlags(), flags)
}

func TestNew(t *testing.T) {
	type args struct {
		githubOwner bool
		githubRepo  bool
	}
	tests := []struct {
		name          string
		args          args
		retErrControl testclient.ReturnValueStr
		wantErr       error
	}{
		{
			name: "Proper CLI flags given",
			args: args{
				githubOwner: true,
				githubRepo:  true,
			},
			wantErr: nil,
		},
		{
			name: "github-owner flag is missing",
			args: args{
				githubOwner: false,
				githubRepo:  true,
			},
			wantErr: errors.New("option --github-owner is required"),
		},
		{
			name: "github-repo flag is missing",
			args: args{
				githubOwner: true,
				githubRepo:  false,
			},
			wantErr: errors.New("option --github-repo is required"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			github.NewGitHubClientFunc = testclient.New
			testclient.ReturnValue = tt.retErrControl
			_, err := github.New(setupCLIContext(tt.args.githubOwner, tt.args.githubRepo))

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("New() got = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
