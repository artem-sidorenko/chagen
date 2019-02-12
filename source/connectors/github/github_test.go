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

	tcli "github.com/artem-sidorenko/chagen/internal/testing/cli"
	"github.com/artem-sidorenko/chagen/source/connectors"
	"github.com/artem-sidorenko/chagen/source/connectors/github"
	"github.com/artem-sidorenko/chagen/source/connectors/github/internal/testclient"
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
				RetRepoServiceGetErr:      true,
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
			cliFlags := map[string]string{}
			if tt.args.githubOwner {
				cliFlags["github-owner"] = "testowner"
			}

			if tt.args.githubRepo {
				cliFlags["github-repo"] = "testrepo"
			}

			ctx := tcli.TestContext(github.CLIFlags(), cliFlags)

			github.NewGitHubClientFunc = testclient.New
			testclient.ReturnValue = tt.retErrControl
			_, err := github.New(ctx)

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("New() got = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
