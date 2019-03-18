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

package data_test

import (
	"reflect"
	"testing"

	"github.com/artem-sidorenko/chagen/data"
	"github.com/artem-sidorenko/chagen/internal/testing/helpers"
	"github.com/artem-sidorenko/chagen/internal/testing/testdata"
)

func TestNewReleases(t *testing.T) {
	tests := []struct {
		name    string
		wantRet data.Releases
	}{
		{
			name: "proper data with all elements",
			wantRet: data.Releases{
				data.Release{
					Release:    "v0.1.2",
					Date:       "20.03.2003",
					ReleaseURL: "https://test.example.com/tags/v0.1.2",
					MRs: data.MRs{
						data.MR{
							ID:         2314,
							Name:       "Test PR title 11",
							URL:        "https://test.example.com/mrs/2314",
							MergedDate: helpers.Time(1048094647),
							Author:     "test-user8",
							AuthorURL:  "https://test.example.com/authors/test-user8",
							Labels:     []string{"no changelog"},
						},
					},
				},
				data.Release{
					Release:    "v0.1.1",
					Date:       "19.03.2003",
					ReleaseURL: "https://test.example.com/tags/v0.1.1",
					Issues: data.Issues{
						data.Issue{
							ID:         1304,
							Name:       "Test issue title 10",
							ClosedDate: helpers.Time(1047993647),
							URL:        "http://test.example.com/issues/1304",
							Labels:     []string{"wontfix"},
						},
					},
					MRs: data.MRs{
						data.MR{
							ID:         2304,
							Name:       "Test PR title 10",
							URL:        "https://test.example.com/mrs/2304",
							MergedDate: helpers.Time(1047994647),
							Author:     "test-user",
							AuthorURL:  "https://test.example.com/authors/test-user",
							Labels:     []string{"bugfix"},
						},
					},
				},
				data.Release{
					Release:    "v0.1.0",
					Date:       "18.03.2003",
					ReleaseURL: "https://test.example.com/tags/v0.1.0",
					Issues: data.Issues{
						data.Issue{
							ID:         1294,
							Name:       "Test issue title 9",
							ClosedDate: helpers.Time(1047893647),
							URL:        "http://test.example.com/issues/1294",
						},
					},
					MRs: data.MRs{
						data.MR{
							ID:         2294,
							Name:       "Test PR title 9",
							URL:        "https://test.example.com/mrs/2294",
							MergedDate: helpers.Time(1047894647),
							Author:     "test-user",
							AuthorURL:  "https://test.example.com/authors/test-user",
							Labels:     []string{"bugfix"},
						},
					},
				},
				data.Release{
					Release:    "v0.0.9",
					Date:       "17.03.2003",
					ReleaseURL: "https://test.example.com/tags/v0.0.9",
					MRs: data.MRs{
						data.MR{
							ID:         2284,
							Name:       "Test PR title 8",
							URL:        "https://test.example.com/mrs/2284",
							MergedDate: helpers.Time(1047794647),
							Author:     "test-user",
							AuthorURL:  "https://test.example.com/authors/test-user",
							Labels:     []string{"invalid"},
						},
					},
				},
				data.Release{
					Release:    "v0.0.8",
					Date:       "16.03.2003",
					ReleaseURL: "https://test.example.com/tags/v0.0.8",
					Issues: data.Issues{
						data.Issue{
							ID:         1274,
							Name:       "Test issue title 7",
							ClosedDate: helpers.Time(1047693647),
							URL:        "http://test.example.com/issues/1274",
							Labels:     []string{"no changelog"},
						},
					},
					MRs: data.MRs{
						data.MR{
							ID:         2274,
							Name:       "Test PR title 7",
							URL:        "https://test.example.com/mrs/2274",
							MergedDate: helpers.Time(1047694647),
							Author:     "test5-user",
							AuthorURL:  "https://test.example.com/authors/test5-user",
							Labels:     []string{"bugfix"},
						},
					},
				},
				data.Release{
					Release:    "v0.0.7",
					Date:       "14.03.2003",
					ReleaseURL: "https://test.example.com/tags/v0.0.7",
					Issues: data.Issues{
						data.Issue{
							ID:         1264,
							Name:       "Test issue title 6",
							ClosedDate: helpers.Time(1047593647),
							URL:        "http://test.example.com/issues/1264",
							Labels:     []string{"invalid"},
						},
					},
					MRs: data.MRs{
						data.MR{
							ID:         2264,
							Name:       "Test PR title 6",
							URL:        "https://test.example.com/mrs/2264",
							MergedDate: helpers.Time(1047594647),
							Author:     "test-user",
							AuthorURL:  "https://test.example.com/authors/test-user",
							Labels:     []string{"enhancement"},
						},
					},
				},
				data.Release{
					Release:    "v0.0.6",
					Date:       "13.03.2003",
					ReleaseURL: "https://test.example.com/tags/v0.0.6",
					MRs: data.MRs{
						data.MR{
							ID:         2254,
							Name:       "Test PR title 5",
							URL:        "https://test.example.com/mrs/2254",
							MergedDate: helpers.Time(1047494647),
							Author:     "test-user",
							AuthorURL:  "https://test.example.com/authors/test-user",
							Labels:     []string{"bugfix"},
						},
					},
				},
				data.Release{
					Release:    "v0.0.5",
					Date:       "12.03.2003",
					ReleaseURL: "https://test.example.com/tags/v0.0.5",
					Issues: data.Issues{
						data.Issue{
							ID:         1244,
							Name:       "Test issue title 4",
							ClosedDate: helpers.Time(1047393647),
							URL:        "http://test.example.com/issues/1244",
						},
					},
				},
				data.Release{
					Release:    "v0.0.4",
					Date:       "11.03.2003",
					ReleaseURL: "https://test.example.com/tags/v0.0.4",
					MRs: data.MRs{
						data.MR{
							ID:         2234,
							Name:       "Test PR title 3",
							URL:        "https://test.example.com/mrs/2234",
							MergedDate: helpers.Time(1047294647),
							Author:     "test-user",
							AuthorURL:  "https://test.example.com/authors/test-user",
							Labels:     []string{"enhancement", "bugfix"},
						},
					},
				},
				data.Release{
					Release:    "v0.0.3",
					Date:       "10.03.2003",
					ReleaseURL: "https://test.example.com/tags/v0.0.3",
					Issues: data.Issues{
						data.Issue{
							ID:         1227,
							Name:       "Test issue title 2",
							ClosedDate: helpers.Time(1047193647),
							URL:        "http://test.example.com/issues/1227",
							Labels:     []string{"enhancement", "bugfix"},
						},
					},
					MRs: data.MRs{
						data.MR{
							ID:         2224,
							Name:       "Test PR title 2",
							URL:        "https://test.example.com/mrs/2224",
							MergedDate: helpers.Time(1047194647),
							Author:     "test-user2",
							AuthorURL:  "https://test.example.com/authors/test-user2",
						},
					},
				},
				data.Release{
					Release:    "v0.0.2",
					Date:       "09.03.2003",
					ReleaseURL: "https://test.example.com/tags/v0.0.2",
					Issues: data.Issues{
						data.Issue{
							ID:         1214,
							Name:       "Test issue title 1",
							ClosedDate: helpers.Time(1047093647),
							URL:        "http://test.example.com/issues/1214",
							Labels:     []string{"enhancement"},
						},
					},
					MRs: data.MRs{
						data.MR{
							ID:         2214,
							Name:       "Test PR title 1",
							URL:        "https://test.example.com/mrs/2214",
							MergedDate: helpers.Time(1047094647),
							Author:     "test-user",
							AuthorURL:  "https://test.example.com/authors/test-user",
							Labels:     []string{"bugfix"},
						},
					},
				},
				data.Release{
					Release:    "v0.0.1",
					Date:       "08.03.2003",
					ReleaseURL: "https://test.example.com/tags/v0.0.1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags := testdata.DataTags()
			issues := testdata.DataIssues()
			mrs := testdata.DataMRs()
			gotRet := data.NewReleases(tags, issues, mrs)
			if len(gotRet) != len(tt.wantRet) {
				t.Errorf("NewReleases() different amount of results. got %#v, want %#v",
					len(gotRet), len(tt.wantRet))
				t.FailNow()
			}
			for i := range gotRet {
				if !reflect.DeepEqual(gotRet[i], tt.wantRet[i]) {
					t.Errorf("\nNewReleases() [%v] = \n got %#v, \n want %#v", i, gotRet[i], tt.wantRet[i])
				}
			}
		})
	}
}
