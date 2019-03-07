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
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/artem-sidorenko/chagen/data"
	"github.com/artem-sidorenko/chagen/datasource/connectors/gitlab"
	"github.com/artem-sidorenko/chagen/datasource/connectors/gitlab/internal/testclient"
	"github.com/artem-sidorenko/chagen/internal/testing/helpers"
)

func TestConnector_Issues(t *testing.T) {
	tests := []struct {
		name          string
		returnValue   testclient.ReturnValueStr
		want          data.Issues
		wantErr       error
		wantMaxIssues []int
	}{
		{
			name: "API returns proper data",
			want: data.Issues{
				data.Issue{
					ID:         1234,
					Name:       "Test issue title 13",
					ClosedDate: time.Unix(2048293647, 0),
					URL:        "https://example.com/issues/1234",
					Labels:     []string{"enhancement"},
				},
				data.Issue{
					ID:         1224,
					Name:       "Test issue title 12",
					ClosedDate: time.Unix(2048193647, 0),
					URL:        "https://example.com/issues/1224",
					Labels:     []string(nil),
				},
				data.Issue{
					ID:         1304,
					Name:       "Test issue title 10",
					ClosedDate: time.Unix(2047993647, 0),
					URL:        "https://example.com/issues/1304",
					Labels:     []string{"wontfix"},
				},
				data.Issue{
					ID:         1294,
					Name:       "Test issue title 9",
					ClosedDate: time.Unix(2047893647, 0),
					URL:        "https://example.com/issues/1294",
					Labels:     []string(nil),
				},
				data.Issue{
					ID:         1274,
					Name:       "Test issue title 7",
					ClosedDate: time.Unix(2047693647, 0),
					URL:        "https://example.com/issues/1274",
					Labels:     []string{"no changelog"},
				},
				data.Issue{
					ID:         1264,
					Name:       "Test issue title 6",
					ClosedDate: time.Unix(2047593647, 0),
					URL:        "https://example.com/issues/1264",
					Labels:     []string{"invalid"},
				},
				data.Issue{
					ID:         1244,
					Name:       "Test issue title 4",
					ClosedDate: time.Unix(2047393647, 0),
					URL:        "https://example.com/issues/1244",
					Labels:     []string(nil),
				},
				data.Issue{
					ID:         1227,
					Name:       "Test issue title 2",
					ClosedDate: time.Unix(2047193647, 0),
					URL:        "https://example.com/issues/1227",
					Labels:     []string{"enhancement", "bugfix"},
				},
				data.Issue{
					ID:         1214,
					Name:       "Test issue title 1",
					ClosedDate: time.Unix(2047093647, 0),
					URL:        "https://example.com/issues/1214",
					Labels:     []string{"enhancement"},
				},
			},
			wantMaxIssues: []int{9},
		},
		{
			name: "ListIssues call fails",
			returnValue: testclient.ReturnValueStr{
				IssuesServiceListProjectIssuesErr: true,
			},
			wantErr: errors.New("can't fetch the issues"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitlab.IssuesPerPage = 5
			c := setupTestConnector(tt.returnValue)
			cerr := make(chan error, 1)

			cgot, _, cmaxissues := c.Issues(context.Background(), cerr)
			gotmaxissues := helpers.GetChannelValuesInt(cmaxissues)

			var got data.Issues
			for t := range cgot {
				got = append(got, t)
			}
			// sort the issues to have the stable order
			sort.Sort(&got)

			// sleep and allow the possible error to be delivered to the channel
			time.Sleep(time.Millisecond * 200)
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
				for i, g := range got {
					fmt.Printf("%#v \n %#v \n\n", g, tt.want[i])
				}
			}

			if err == nil { // compare the processed Issues only in non-error situation
				if !reflect.DeepEqual(gotmaxissues, tt.wantMaxIssues) {
					t.Errorf("Connector.Issues() maxissues = %v, want %v", gotmaxissues, tt.wantMaxIssues)
				}
			}
		})
	}
}