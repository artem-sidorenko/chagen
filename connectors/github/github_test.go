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
	"github.com/artem-sidorenko/chagen/connectors/github"
	"github.com/artem-sidorenko/chagen/connectors/github/internal/testclient"
	"github.com/artem-sidorenko/chagen/data"
	tcli "github.com/artem-sidorenko/chagen/internal/testing/cli"

	"github.com/urfave/cli"
)

func setupTestConnector(
	retErrControl testclient.RetErr, newTagUseReleaseURL bool,
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
	testclient.RetErrControl = retErrControl

	c, _ := github.New(ctx)

	return c
}

func TestGetTags(t *testing.T) {
	tests := []struct {
		name          string
		retErrControl testclient.RetErr
		want          data.Tags
		wantErr       error
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
			retErrControl: testclient.RetErr{
				RetRepoServiceListTagsErr: true,
			},
			wantErr: errors.New("GitHub query 'GetTags' failed: Can't fetch the tags"),
		},
		{
			name: "GetCommit call fails",
			retErrControl: testclient.RetErr{
				RetRepoServiceGetCommitsErr: true,
			},
			wantErr: errors.New("GitHub query 'GetTags' failed: Can't fetch the commit"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestConnector(tt.retErrControl, false)

			got, err := c.GetTags()

			if err != nil && err.Error() != tt.wantErr.Error() {
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
		name          string
		retErrControl testclient.RetErr
		want          data.Issues
		wantErr       error
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
			retErrControl: testclient.RetErr{
				RetIssueServiceListByRepoErr: true,
			},
			wantErr: errors.New("GitHub query 'GetIssues' failed: Can't fetch the issues"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestConnector(tt.retErrControl, false)

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

func TestGetMRs(t *testing.T) {
	tests := []struct {
		name          string
		retErrControl testclient.RetErr
		want          data.MRs
		wantErr       error
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
			retErrControl: testclient.RetErr{
				RetPullRequestsListErr: true,
			},
			wantErr: errors.New("GitHub query 'GetMRs' failed: Can't fetch the PRs"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestConnector(tt.retErrControl, false)

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
			c := setupTestConnector(testclient.RetErr{}, tt.fields.NewTagUseReleaseURL)

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
		name    string
		args    args
		wantErr error
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
			wantErr: errors.New("Option --github-owner is required"),
		},
		{
			name: "github-repo flag is missing",
			args: args{
				githubOwner: true,
				githubRepo:  false,
			},
			wantErr: errors.New("Option --github-repo is required"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := github.New(setupCLIContext(tt.args.githubOwner, tt.args.githubRepo))

			if (err != nil && tt.wantErr == nil) ||
				(err == nil && tt.wantErr != nil) ||
				((err != nil && tt.wantErr != nil) && (err.Error() != tt.wantErr.Error())) {

				t.Errorf("New() got = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
