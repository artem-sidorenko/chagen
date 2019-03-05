/*
   Copyright 2019 Artem Sidorenko <artem@posteo.de>

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

package gitlab_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	tcli "github.com/artem-sidorenko/chagen/internal/testing/cli"
	"github.com/artem-sidorenko/chagen/source/connectors"
	"github.com/artem-sidorenko/chagen/source/connectors/gitlab"
	"github.com/artem-sidorenko/chagen/source/connectors/gitlab/internal/testclient"
)

func setupTestConnector(
	returnValue testclient.ReturnValueStr,
) connectors.Connector {

	gitlab.NewClient = testclient.New
	cliFlags := map[string]string{
		"gitlab-owner": "testowner",
		"gitlab-repo":  "testrepo",
	}

	ctx := tcli.TestContext(gitlab.CLIFlags(), cliFlags)

	testclient.ReturnValue = returnValue

	c, err := gitlab.New(ctx)
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
				ProjectsServiceGetProjectRespCode: 200,
			},
			want: true,
		},
		{
			name: "API returns 404",
			returnValue: testclient.ReturnValueStr{
				ProjectsServiceGetProjectRespCode: 404,
				ProjectsServiceGetProjectErr:      true,
			},
			want: false,
		},
		{
			name: "API returns unhandled error code 500",
			returnValue: testclient.ReturnValueStr{
				ProjectsServiceGetProjectRespCode: 500,
			},
			wantErr: errors.New("GitLab query 'RepositoryExists' failed: unhandled HTTP response code 500"),
		},
		{
			name: "Get returns an error",
			returnValue: testclient.ReturnValueStr{
				ProjectsServiceGetProjectErr: true,
			},
			wantErr: errors.New("GitLab query 'RepositoryExists' failed: can't fetch the repo data"),
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestConnector(tt.returnValue)

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
		gitlabOwner bool
		gitlabRepo  bool
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
				gitlabOwner: true,
				gitlabRepo:  true,
			},
			wantErr: nil,
		},
		{
			name: "gitlab-owner flag is missing",
			args: args{
				gitlabOwner: false,
				gitlabRepo:  true,
			},
			wantErr: errors.New("option --gitlab-owner is required"),
		},
		{
			name: "gitlab-repo flag is missing",
			args: args{
				gitlabOwner: true,
				gitlabRepo:  false,
			},
			wantErr: errors.New("option --gitlab-repo is required"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cliFlags := map[string]string{}
			if tt.args.gitlabOwner {
				cliFlags["gitlab-owner"] = "testowner"
			}

			if tt.args.gitlabRepo {
				cliFlags["gitlab-repo"] = "testrepo"
			}

			ctx := tcli.TestContext(gitlab.CLIFlags(), cliFlags)

			gitlab.NewClient = testclient.New
			testclient.ReturnValue = tt.retErrControl
			_, err := gitlab.New(ctx)

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("New() got = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
