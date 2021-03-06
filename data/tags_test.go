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
	"regexp"
	"sort"
	"testing"
	"time"

	"github.com/artem-sidorenko/chagen/data"
)

func TestTags_Sort(t *testing.T) {
	tests := []struct {
		name string
		t    *data.Tags
		want *data.Tags
	}{
		{
			name: "Tags are already sorted in the wrong order",
			t: &data.Tags{
				{
					Name: "v0.0.1",
					Date: time.Unix(1047483647, 0),
				},
				{
					Name: "v0.0.2",
					Date: time.Unix(1147483647, 0),
				},
				{
					Name: "v0.0.3",
					Date: time.Unix(1247483647, 0),
				},
			},
			want: &data.Tags{
				{
					Name: "v0.0.3",
					Date: time.Unix(1247483647, 0),
				},
				{
					Name: "v0.0.2",
					Date: time.Unix(1147483647, 0),
				},
				{
					Name: "v0.0.1",
					Date: time.Unix(1047483647, 0),
				},
			},
		},

		{
			name: "Tags are not sorted",
			t: &data.Tags{
				{
					Name: "v0.0.2",
					Date: time.Unix(1147483647, 0),
				},
				{
					Name: "v0.0.1",
					Date: time.Unix(1047483647, 0),
				},
				{
					Name: "v0.0.3",
					Date: time.Unix(1247483647, 0),
				},
			},
			want: &data.Tags{
				{
					Name: "v0.0.3",
					Date: time.Unix(1247483647, 0),
				},
				{
					Name: "v0.0.2",
					Date: time.Unix(1147483647, 0),
				},
				{
					Name: "v0.0.1",
					Date: time.Unix(1047483647, 0),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Sort(tt.t)

			if !reflect.DeepEqual(tt.t, tt.want) {
				t.Errorf("sort.Sort(Tags), got %v, want %v", tt.t, tt.want)
			}
		})
	}
}

func TestFilterTags(t *testing.T) {
	type args struct {
		ts     data.Tags
		regexp *regexp.Regexp
	}
	tests := []struct {
		name string
		args args
		want data.Tags
	}{
		{
			name: "Filtering of tags",
			args: args{
				regexp: regexp.MustCompile(`^v\d+\.\d+\.\d+$`),
				ts: data.Tags{
					{
						Name: "asdasdv1.16.55",
						Date: time.Unix(1147493347, 0),
					},
					{
						Name: "v1.16.55",
						Date: time.Unix(1147483347, 0),
					},
					{
						Name: "v0.0.2",
						Date: time.Unix(1147483647, 0),
					},
					{
						Name: "testing",
						Date: time.Unix(1047483647, 0),
					},
					{
						Name: "v0.0.3",
						Date: time.Unix(1247483647, 0),
					},
				},
			},
			want: data.Tags{
				{
					Name: "v1.16.55",
					Date: time.Unix(1147483347, 0),
				},
				{
					Name: "v0.0.2",
					Date: time.Unix(1147483647, 0),
				},
				{
					Name: "v0.0.3",
					Date: time.Unix(1247483647, 0),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := data.FilterTags(tt.args.ts, tt.args.regexp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterIssuesByName() = %v, want %v", got, tt.want)
			}
		})
	}
}
