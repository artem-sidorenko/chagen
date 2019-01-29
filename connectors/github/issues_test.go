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
					ID:         1214,
					Name:       "Test issue title 1",
					ClosedDate: time.Unix(2047093647, 0),
					URL:        "http://example.com/issues/1214",
					Labels:     []string{"enhancement"},
				},
				data.Issue{
					ID:         1227,
					Name:       "Test issue title 2",
					ClosedDate: time.Unix(2047193647, 0),
					URL:        "http://example.com/issues/1227",
					Labels:     []string{"enhancement", "bugfix"},
				},
				data.Issue{
					ID:         1244,
					Name:       "Test issue title 4",
					ClosedDate: time.Unix(2047393647, 0),
					URL:        "http://example.com/issues/1244",
					Labels:     []string(nil),
				},
				data.Issue{
					ID:         1264,
					Name:       "Test issue title 6",
					ClosedDate: time.Unix(2047593647, 0),
					URL:        "http://example.com/issues/1264",
					Labels:     []string{"invalid"},
				},
				data.Issue{
					ID:         1274,
					Name:       "Test issue title 7",
					ClosedDate: time.Unix(2047693647, 0),
					URL:        "http://example.com/issues/1274",
					Labels:     []string{"no changelog"},
				},
				data.Issue{
					ID:         1294,
					Name:       "Test issue title 9",
					ClosedDate: time.Unix(2047893647, 0),
					URL:        "http://example.com/issues/1294",
					Labels:     []string(nil),
				},
				data.Issue{
					ID:         1304,
					Name:       "Test issue title 10",
					ClosedDate: time.Unix(2047993647, 0),
					URL:        "http://example.com/issues/1304",
					Labels:     []string{"wontfix"},
				},
				data.Issue{
					ID:         1224,
					Name:       "Test issue title 12",
					ClosedDate: time.Unix(2048193647, 0),
					URL:        "http://example.com/issues/1224",
					Labels:     []string(nil),
				},
				data.Issue{
					ID:         1234,
					Name:       "Test issue title 13",
					ClosedDate: time.Unix(2048293647, 0),
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

			cgot, _, _ := c.Issues(context.Background(), cerr)
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
