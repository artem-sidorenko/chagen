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
	"sort"
	"testing"
	"time"

	"github.com/artem-sidorenko/chagen/data"
)

func TestIssues_Sort(t *testing.T) {
	tests := []struct {
		name string
		is   *data.Issues
		want *data.Issues
	}{
		{
			name: "Issues are already sorted in the wrong order",
			is: &data.Issues{
				{
					Name:       "Issue 1",
					ClosedDate: time.Unix(1047483647, 0),
				},
				{
					Name:       "Issue 2",
					ClosedDate: time.Unix(1247483647, 0),
				},
				{
					Name:       "Issue 3",
					ClosedDate: time.Unix(1347483647, 0),
				},
			},
			want: &data.Issues{
				{
					Name:       "Issue 3",
					ClosedDate: time.Unix(1347483647, 0),
				},
				{
					Name:       "Issue 2",
					ClosedDate: time.Unix(1247483647, 0),
				},
				{
					Name:       "Issue 1",
					ClosedDate: time.Unix(1047483647, 0),
				},
			},
		},
		{
			name: "Issues are not sorted",
			is: &data.Issues{
				{
					Name:       "Issue 2",
					ClosedDate: time.Unix(1247483647, 0),
				},
				{
					Name:       "Issue 1",
					ClosedDate: time.Unix(1047483647, 0),
				},
				{
					Name:       "Issue 3",
					ClosedDate: time.Unix(1347483647, 0),
				},
			},
			want: &data.Issues{
				{
					Name:       "Issue 3",
					ClosedDate: time.Unix(1347483647, 0),
				},
				{
					Name:       "Issue 2",
					ClosedDate: time.Unix(1247483647, 0),
				},
				{
					Name:       "Issue 1",
					ClosedDate: time.Unix(1047483647, 0),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Sort(tt.is)

			if !reflect.DeepEqual(tt.is, tt.want) {
				t.Errorf("sort.Sort(Issues), got%v, want %v", tt.is, tt.want)
			}
		})
	}
}

func TestIssues_Filter(t *testing.T) {
	type args struct {
		fromDate time.Time
		toDate   time.Time
	}
	tests := []struct {
		name    string
		is      data.Issues
		args    args
		wantRet data.Issues
	}{
		{
			name: "Filtering of issues",
			is: data.Issues{
				{
					Name:       "Issue 1",
					ClosedDate: time.Unix(1047483647, 0),
				},
				{
					Name:       "Issue 2",
					ClosedDate: time.Unix(1247483647, 0),
				},
				{
					Name:       "Issue 3",
					ClosedDate: time.Unix(1347483647, 0),
				},
			},
			args: args{
				fromDate: time.Unix(1057483647, 0),
				toDate:   time.Unix(1337483647, 0),
			},
			wantRet: data.Issues{
				{
					Name:       "Issue 2",
					ClosedDate: time.Unix(1247483647, 0),
				},
			},
		},
		{
			name: "Filtering of same day issue",
			is: data.Issues{
				{
					Name:       "Issue 1",
					ClosedDate: time.Unix(1047483647, 0),
				},
				{
					Name:       "Issue 2",
					ClosedDate: time.Unix(1247483647, 0),
				},
				{
					Name:       "Issue 3",
					ClosedDate: time.Unix(1347483647, 0),
				},
			},
			args: args{
				fromDate: time.Unix(1057483647, 0),
				toDate:   time.Unix(1247483647, 0),
			},
			wantRet: data.Issues{
				{
					Name:       "Issue 2",
					ClosedDate: time.Unix(1247483647, 0),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet := data.FilterIssues(tt.is, tt.args.fromDate, tt.args.toDate)
			if !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("FilterIssues(), got %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}
