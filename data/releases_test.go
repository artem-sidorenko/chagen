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

	"github.com/artem-sidorenko/chagen/connectors"
	"github.com/artem-sidorenko/chagen/data"
)

func TestNewReleases(t *testing.T) {
	type args struct {
		tags   connectors.Tags
		issues connectors.Issues
		mrs    connectors.MRs
	}
	tests := []struct {
		name    string
		args    args
		wantRet data.Releases
	}{
		{
			name: "proper data with all elements",
			args: args{
				tags: connectors.Tags{
					connectors.Tag{
						Name:   "v0.0.1",
						Date:   time.Unix(1047483647, 0),
						Commit: "b6a735dcb420a82865abe8c194900e59f6af9dea",
						URL:    "https://example.com/tags/v0.0.1",
					},
				},
				issues: connectors.Issues{
					connectors.Issue{
						Name:       "Issue number one",
						ClosedDate: time.Unix(1047482647, 0),
						ID:         1,
						URL:        "https://example.com/issues/1",
					},
				},
				mrs: connectors.MRs{
					connectors.MR{
						Name:       "Test merge request",
						ID:         2,
						MergedDate: time.Unix(1047480647, 0),
						Author:     "testauthor",
						AuthorURL:  "https://example.com/authors/testauthor",
						URL:        "https://example.com/mrs/2",
					},
				},
			},
			wantRet: data.Releases{
				data.Release{
					Release:    "v0.0.1",
					Date:       "12.03.2003",
					ReleaseURL: "https://example.com/tags/v0.0.1",
					Issues: data.Issues{
						data.Issue{
							ID:   1,
							Name: "Issue number one",
							URL:  "https://example.com/issues/1",
						},
					},
					MRs: data.MRs{
						data.MR{
							ID:        2,
							Name:      "Test merge request",
							URL:       "https://example.com/mrs/2",
							Author:    "testauthor",
							AuthorURL: "https://example.com/authors/testauthor",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet := data.NewReleases(tt.args.tags, tt.args.issues, tt.args.mrs)
			if !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("NewReleases() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}
