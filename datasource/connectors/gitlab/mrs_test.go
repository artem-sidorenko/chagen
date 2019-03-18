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
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/artem-sidorenko/chagen/data"
	"github.com/artem-sidorenko/chagen/datasource/connectors/gitlab"
	"github.com/artem-sidorenko/chagen/datasource/connectors/gitlab/internal/testclient"
	"github.com/artem-sidorenko/chagen/internal/testing/helpers"
)

func TestConnector_MRs(t *testing.T) {
	tests := []struct {
		name        string
		returnValue testclient.ReturnValueStr
		want        data.MRs
		wantErr     error
		wantMaxMRs  []int
	}{
		{
			name: "API returns proper data",
			want: data.MRs{
				data.MR{
					ID:         2344,
					Name:       "Test PR title 14",
					URL:        "https://example.com/pulls/2344",
					Author:     "te77st-user",
					AuthorURL:  "https://gitlab.com/te77st-user",
					MergedDate: helpers.Time(1048394647),
					Labels:     []string{"bugfix"},
				},
				data.MR{
					ID:         2334,
					Name:       "Test PR title 13",
					URL:        "https://example.com/pulls/2334",
					Author:     "test-user",
					AuthorURL:  "https://gitlab.com/test-user",
					MergedDate: helpers.Time(1048294647),
					Labels:     []string{"bugfix"},
				},
				data.MR{
					ID:         2314,
					Name:       "Test PR title 11",
					URL:        "https://example.com/pulls/2314",
					Author:     "test-user8",
					AuthorURL:  "https://gitlab.com/test-user8",
					MergedDate: helpers.Time(1048094647),
					Labels:     []string{"no changelog"},
				},
				data.MR{
					ID:         2304,
					Name:       "Test PR title 10",
					URL:        "https://example.com/pulls/2304",
					Author:     "test-user",
					AuthorURL:  "https://gitlab.com/test-user",
					MergedDate: helpers.Time(1047994647),
					Labels:     []string{"bugfix"},
				},
				data.MR{
					ID:         2294,
					Name:       "Test PR title 9",
					URL:        "https://example.com/pulls/2294",
					Author:     "test-user",
					AuthorURL:  "https://gitlab.com/test-user",
					MergedDate: helpers.Time(1047894647),
					Labels:     []string{"bugfix"},
				},
				data.MR{
					ID:         2284,
					Name:       "Test PR title 8",
					URL:        "https://example.com/pulls/2284",
					Author:     "test-user",
					AuthorURL:  "https://gitlab.com/test-user",
					MergedDate: helpers.Time(1047794647),
					Labels:     []string{"invalid"},
				},
				data.MR{
					ID:         2274,
					Name:       "Test PR title 7",
					URL:        "https://example.com/pulls/2274",
					Author:     "test5-user",
					AuthorURL:  "https://gitlab.com/test5-user",
					MergedDate: helpers.Time(1047694647),
					Labels:     []string{"bugfix"},
				},
				data.MR{
					ID:         2264,
					Name:       "Test PR title 6",
					URL:        "https://example.com/pulls/2264",
					Author:     "test-user",
					AuthorURL:  "https://gitlab.com/test-user",
					MergedDate: helpers.Time(1047594647),
					Labels:     []string{"enhancement"},
				},
				data.MR{
					ID:         2254,
					Name:       "Test PR title 5",
					URL:        "https://example.com/pulls/2254",
					Author:     "test-user",
					AuthorURL:  "https://gitlab.com/test-user",
					MergedDate: helpers.Time(1047494647),
					Labels:     []string{"bugfix"},
				},
				data.MR{
					ID:         2234,
					Name:       "Test PR title 3",
					URL:        "https://example.com/pulls/2234",
					Author:     "test-user",
					AuthorURL:  "https://gitlab.com/test-user",
					MergedDate: helpers.Time(1047294647),
					Labels:     []string{"enhancement", "bugfix"},
				},
				data.MR{
					ID:         2224,
					Name:       "Test PR title 2",
					URL:        "https://example.com/pulls/2224",
					Author:     "test-user2",
					AuthorURL:  "https://gitlab.com/test-user2",
					MergedDate: helpers.Time(1047194647),
					Labels:     []string(nil),
				},
				data.MR{
					ID:         2214,
					Name:       "Test PR title 1",
					URL:        "https://example.com/pulls/2214",
					Author:     "test-user",
					AuthorURL:  "https://gitlab.com/test-user",
					MergedDate: helpers.Time(1047094647),
					Labels:     []string{"bugfix"},
				},
			},
			wantMaxMRs: []int{12},
		},
		{
			name: "ListProjectMergeRequests call fails",
			returnValue: testclient.ReturnValueStr{
				MergeRequestsServiceListProjectMergeRequestsErr: true,
			},
			wantErr: errors.New("can't fetch the MRs"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitlab.MRsPerPage = 5
			c := setupTestConnector(tt.returnValue)
			cerr := make(chan error, 1)

			cgot, _, cmaxmrs := c.MRs(context.Background(), cerr)
			gotmaxmrs := helpers.GetChannelValuesInt(cmaxmrs)

			var got data.MRs
			for t := range cgot {
				got = append(got, t)
			}
			// sort the mrs to have the stable order
			sort.Sort(&got)

			// sleep and allow the possible error to be delivered to the channel
			time.Sleep(time.Millisecond * 200)
			var err error
			select {
			case err = <-cerr:
			default:
			}

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Connector.MRs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connector.MRs() = %+v,\n want %+v", got, tt.want)
			}

			if err == nil { // compare the processed MRs only in non-error situation
				if !reflect.DeepEqual(gotmaxmrs, tt.wantMaxMRs) {
					t.Errorf("Connector.MRs() maxmrs = %v, want %v", gotmaxmrs, tt.wantMaxMRs)
				}
			}
		})
	}
}
