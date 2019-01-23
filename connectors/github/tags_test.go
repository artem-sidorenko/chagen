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

package github_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/artem-sidorenko/chagen/connectors/github/internal/testclient"
	"github.com/artem-sidorenko/chagen/data"
)

func TestConnector_Tags(t *testing.T) {
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
			wantErr: errors.New("can't fetch the tags"),
		},
		{
			name: "GetCommit call fails",
			returnValue: testclient.ReturnValueStr{
				RetRepoServiceGetCommitsErr: true,
			},
			wantErr: errors.New("can't fetch the commit"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestConnector(tt.returnValue, false)
			cerr := make(chan error, 1)

			cgot := c.Tags(context.Background(), cerr, nil)
			var got data.Tags
			for t := range cgot {
				got = append(got, t)
			}

			var err error
			select {
			case err = <-cerr:
			default:
			}

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Connector.Tags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connector.Tags() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
