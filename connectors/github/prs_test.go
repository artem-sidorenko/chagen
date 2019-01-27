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

func TestConnector_MRs(t *testing.T) {
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
					ID:         2214,
					Name:       "Test PR title 1",
					URL:        "https://example.com/pulls/2214",
					Author:     "test-user",
					AuthorURL:  "https://example.com/users/test-user",
					MergedDate: time.Unix(2047094647, 0),
					Labels:     []string{"bugfix"},
				},
				data.MR{
					ID:         2224,
					Name:       "Test PR title 2",
					URL:        "https://example.com/pulls/2224",
					Author:     "test-user2",
					AuthorURL:  "https://example.com/users/test-user2",
					MergedDate: time.Unix(2047194647, 0),
					Labels:     []string(nil),
				},
				data.MR{
					ID:         2234,
					Name:       "Test PR title 3",
					URL:        "https://example.com/pulls/2234",
					Author:     "test-user",
					AuthorURL:  "https://example.com/users/test-user",
					MergedDate: time.Unix(2047294647, 0),
					Labels:     []string{"enhancement", "bugfix"},
				},
				data.MR{
					ID:         2254,
					Name:       "Test PR title 5",
					URL:        "https://example.com/pulls/2254",
					Author:     "test-user",
					AuthorURL:  "https://example.com/users/test-user",
					MergedDate: time.Unix(2047494647, 0),
					Labels:     []string{"bugfix"},
				},
				data.MR{
					ID:         2264,
					Name:       "Test PR title 6",
					URL:        "https://example.com/pulls/2264",
					Author:     "test-user",
					AuthorURL:  "https://example.com/users/test-user",
					MergedDate: time.Unix(2047594647, 0),
					Labels:     []string{"enhancement"},
				},
				data.MR{
					ID:         2274,
					Name:       "Test PR title 7",
					URL:        "https://example.com/pulls/2274",
					Author:     "test5-user",
					AuthorURL:  "https://example.com/users/test5-user",
					MergedDate: time.Unix(2047694647, 0),
					Labels:     []string{"bugfix"},
				},
				data.MR{
					ID:         2284,
					Name:       "Test PR title 8",
					URL:        "https://example.com/pulls/2284",
					Author:     "test-user",
					AuthorURL:  "https://example.com/users/test-user",
					MergedDate: time.Unix(2047794647, 0),
					Labels:     []string{"invalid"},
				},
				data.MR{
					ID:         2294,
					Name:       "Test PR title 9",
					URL:        "https://example.com/pulls/2294",
					Author:     "test-user",
					AuthorURL:  "https://example.com/users/test-user",
					MergedDate: time.Unix(2047894647, 0),
					Labels:     []string{"bugfix"},
				},
				data.MR{
					ID:         2304,
					Name:       "Test PR title 10",
					URL:        "https://example.com/pulls/2304",
					Author:     "test-user",
					AuthorURL:  "https://example.com/users/test-user",
					MergedDate: time.Unix(2047994647, 0),
					Labels:     []string{"bugfix"},
				},
				data.MR{
					ID:         2314,
					Name:       "Test PR title 11",
					URL:        "https://example.com/pulls/2314",
					Author:     "test-user8",
					AuthorURL:  "https://example.com/users/test-user8",
					MergedDate: time.Unix(2048094647, 0),
					Labels:     []string{"no changelog"},
				},
				data.MR{
					ID:         2334,
					Name:       "Test PR title 13",
					URL:        "https://example.com/pulls/2334",
					Author:     "test-user",
					AuthorURL:  "https://example.com/users/test-user",
					MergedDate: time.Unix(2048294647, 0),
					Labels:     []string{"bugfix"},
				},
				data.MR{
					ID:         2344,
					Name:       "Test PR title 14",
					URL:        "https://example.com/pulls/2344",
					Author:     "te77st-user",
					AuthorURL:  "https://example.com/users/te77st-user",
					MergedDate: time.Unix(2048394647, 0),
					Labels:     []string{"bugfix"},
				},
			},
		},
		{
			name: "ListPRs call fails",
			returnValue: testclient.ReturnValueStr{
				RetPullRequestsListErr: true,
			},
			wantErr: errors.New("can't fetch the PRs"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestConnector(tt.returnValue, false)
			cerr := make(chan error, 1)

			cgot, _ := c.MRs(context.Background(), cerr)
			var got data.MRs
			for t := range cgot {
				got = append(got, t)
			}

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
				t.Errorf("Connector.MRs() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
