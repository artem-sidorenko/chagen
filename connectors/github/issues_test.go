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

func TestConnector_Issues(t *testing.T) {
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
			wantErr: errors.New("can't fetch the issues"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestConnector(tt.returnValue, false)
			cerr := make(chan error, 1)

			cgot := c.Issues(context.Background(), cerr, nil)
			var got data.Issues
			for t := range cgot {
				got = append(got, t)
			}

			var err error
			select {
			case err = <-cerr:
			default:
			}

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Connector.Issues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connector.Issues() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
