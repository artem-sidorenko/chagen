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
)

func TestNewReleases(t *testing.T) {
	type args struct {
		tags   data.Tags
		issues data.Issues
		mrs    data.MRs
	}
	tests := []struct {
		name    string
		args    args
		wantRet data.Releases
	}{
		{
			name: "proper data with all elements",
			args: args{
				tags: data.Tags{
					data.Tag{
						Name:   "v0.0.1",
						Date:   time.Unix(1047483647, 0),
						Commit: "b6a735dcb420a82865abe8c194900e59f6af9dea",
						URL:    "https://example.com/tags/v0.0.1",
					},
				},
				issues: data.Issues{
					data.Issue{
						Name:       "Issue number one",
						ClosedDate: time.Unix(1047482647, 0),
						ID:         1,
						URL:        "https://example.com/issues/1",
					},
				},
				mrs: data.MRs{
					data.MR{
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
					Release:       "v0.0.1",
					DateFormatted: "12.03.2003",
					Date:          time.Unix(1047483647, 0),
					ReleaseURL:    "https://example.com/tags/v0.0.1",
					Issues: data.Issues{
						data.Issue{
							ID:         1,
							Name:       "Issue number one",
							ClosedDate: time.Unix(1047482647, 0),
							URL:        "https://example.com/issues/1",
						},
					},
					MRs: data.MRs{
						data.MR{
							ID:         2,
							Name:       "Test merge request",
							URL:        "https://example.com/mrs/2",
							MergedDate: time.Unix(1047480647, 0),
							Author:     "testauthor",
							AuthorURL:  "https://example.com/authors/testauthor",
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

func TestReleases_Sort(t *testing.T) {
	tests := []struct {
		name string
		r    *data.Releases
		want *data.Releases
	}{
		{
			name: "Releases are already sorted",
			r: &data.Releases{
				{
					Release: "v0.0.1",
					Date:    time.Unix(1047483647, 0),
				},
				{
					Release: "v0.0.2",
					Date:    time.Unix(1247483647, 0),
				},
				{
					Release: "v0.0.3",
					Date:    time.Unix(1347483647, 0),
				},
			},
			want: &data.Releases{
				{
					Release: "v0.0.3",
					Date:    time.Unix(1347483647, 0),
				},
				{
					Release: "v0.0.2",
					Date:    time.Unix(1247483647, 0),
				},
				{
					Release: "v0.0.1",
					Date:    time.Unix(1047483647, 0),
				},
			},
		},
		{
			name: "Releases are not sorted",
			r: &data.Releases{
				{
					Release: "v0.0.2",
					Date:    time.Unix(1247483647, 0),
				},
				{
					Release: "v0.0.1",
					Date:    time.Unix(1047483647, 0),
				},
				{
					Release: "v0.0.3",
					Date:    time.Unix(1347483647, 0),
				},
			},
			want: &data.Releases{
				{
					Release: "v0.0.3",
					Date:    time.Unix(1347483647, 0),
				},
				{
					Release: "v0.0.2",
					Date:    time.Unix(1247483647, 0),
				},
				{
					Release: "v0.0.1",
					Date:    time.Unix(1047483647, 0),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Sort()

			if !reflect.DeepEqual(tt.r, tt.want) {
				t.Errorf("Releases.Sort(), Releases = %v, want %v", tt.r, tt.want)
			}
		})
	}
}
