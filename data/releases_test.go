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
	"time"

	"github.com/artem-sidorenko/chagen/data"
	"github.com/artem-sidorenko/chagen/internal/testing/testconnector"
	"github.com/artem-sidorenko/chagen/internal/testing/testconnector/testdata"
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
					Release:    "v0.0.3",
					Date:       "13.07.2009",
					ReleaseURL: "https://test.example.com/tags/v0.0.3",
					Issues: data.Issues{
						data.Issue{
							ID:         2,
							Name:       "Issue 2",
							ClosedDate: time.Unix(1247483647, 0),
							URL:        "http://test.example.com/issues/2",
						},
					},
					MRs: data.MRs{
						data.MR{
							ID:         2,
							Name:       "MR 2",
							URL:        "https://test.example.com/mrs/2",
							MergedDate: time.Unix(1247483647, 0),
							Author:     "testauthor",
							AuthorURL:  "https://test.example.com/authors/testauthor",
						},
					},
				},
				data.Release{
					Release:    "v0.0.2",
					Date:       "13.05.2006",
					ReleaseURL: "https://test.example.com/tags/v0.0.2",
					MRs: data.MRs{
						data.MR{
							ID:         3,
							Name:       "MR 3",
							URL:        "https://test.example.com/mrs/3",
							MergedDate: time.Unix(1057483647, 0),
							Author:     "testauthor",
							AuthorURL:  "https://test.example.com/authors/testauthor",
						},
					},
				},
				data.Release{
					Release:    "v0.0.1",
					Date:       "12.03.2003",
					ReleaseURL: "https://test.example.com/tags/v0.0.1",
					Issues: data.Issues{
						data.Issue{
							ID:         1,
							Name:       "Issue 1",
							ClosedDate: time.Unix(1047483647, 0),
							URL:        "http://test.example.com/issues/1",
						},
					},
					MRs: data.MRs{
						data.MR{
							ID:         1,
							Name:       "MR 1",
							URL:        "https://test.example.com/mrs/1",
							MergedDate: time.Unix(1047483647, 0),
							Author:     "testauthor",
							AuthorURL:  "https://test.example.com/authors/testauthor",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := &testconnector.Connector{}
			tags := testdata.Tags()
			issues := testdata.Issues()
			mrs, _ := conn.GetMRs()
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
