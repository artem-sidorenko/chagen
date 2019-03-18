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
	"github.com/artem-sidorenko/chagen/internal/testing/helpers"
)

func TestMRs_Sort(t *testing.T) {
	tests := []struct {
		name string
		m    *data.MRs
		want *data.MRs
	}{
		{
			name: "MRs are already sorted in the wrong order",
			m: &data.MRs{
				{
					Name:       "MR 1",
					MergedDate: helpers.Time(1047483647),
				},
				{
					Name:       "MR 2",
					MergedDate: helpers.Time(1247483647),
				},
				{
					Name:       "MR 3",
					MergedDate: helpers.Time(1347483647),
				},
			},
			want: &data.MRs{
				{
					Name:       "MR 3",
					MergedDate: helpers.Time(1347483647),
				},
				{
					Name:       "MR 2",
					MergedDate: helpers.Time(1247483647),
				},
				{
					Name:       "MR 1",
					MergedDate: helpers.Time(1047483647),
				},
			},
		},
		{
			name: "MRs are not sorted",
			m: &data.MRs{
				{
					Name:       "MR 2",
					MergedDate: helpers.Time(1247483647),
				},
				{
					Name:       "MR 1",
					MergedDate: helpers.Time(1047483647),
				},
				{
					Name:       "MR 3",
					MergedDate: helpers.Time(1347483647),
				},
			},
			want: &data.MRs{
				{
					Name:       "MR 3",
					MergedDate: helpers.Time(1347483647),
				},
				{
					Name:       "MR 2",
					MergedDate: helpers.Time(1247483647),
				},
				{
					Name:       "MR 1",
					MergedDate: helpers.Time(1047483647),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Sort(tt.m)

			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("sort.Sort(MRs), got %v, want %v", tt.m, tt.want)
			}
		})
	}
}

func TestMRs_Filter(t *testing.T) {
	type args struct {
		fromDate time.Time
		toDate   time.Time
	}
	tests := []struct {
		name    string
		m       data.MRs
		args    args
		wantRet data.MRs
	}{
		{
			name: "Filtering of MRs",
			m: data.MRs{
				{
					Name:       "MR 1",
					MergedDate: helpers.Time(1047483647),
				},
				{
					Name:       "MR 2",
					MergedDate: helpers.Time(1247483647),
				},
				{
					Name:       "MR 3",
					MergedDate: helpers.Time(1347483647),
				},
			},
			args: args{
				fromDate: helpers.Time(1057483647),
				toDate:   helpers.Time(1337483647),
			},
			wantRet: data.MRs{
				{
					Name:       "MR 2",
					MergedDate: helpers.Time(1247483647),
				},
			},
		},
		{
			name: "Filtering of same day MR",
			m: data.MRs{
				{
					Name:       "MR 1",
					MergedDate: helpers.Time(1047483647),
				},
				{
					Name:       "MR 2",
					MergedDate: helpers.Time(1247483647),
				},
				{
					Name:       "MR 3",
					MergedDate: helpers.Time(1347483647),
				},
			},
			args: args{
				fromDate: helpers.Time(1057483647),
				toDate:   helpers.Time(1247483647),
			},
			wantRet: data.MRs{
				{
					Name:       "MR 2",
					MergedDate: helpers.Time(1247483647),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet := data.FilterMRs(tt.m, tt.args.fromDate, tt.args.toDate)
			if !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("FilterMRs(), got %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestFilterMRsByLabel(t *testing.T) {
	type args struct {
		m             data.MRs
		withoutLabels []string
	}
	tests := []struct {
		name string
		args args
		want data.MRs
	}{
		{
			name: "MRs without any matching exclude labels",
			args: args{
				m: data.MRs{
					{
						Name:       "MR 1",
						MergedDate: helpers.Time(1047483647),
						Labels:     []string{"bugfix"},
					},
					{
						Name:       "MR 2",
						MergedDate: helpers.Time(1247483647),
					},
					{
						Name:       "MR 3",
						MergedDate: helpers.Time(1347483647),
						Labels:     []string{"enhancement"},
					},
				},
				withoutLabels: []string{"no changelog", "wontfix"},
			},
			want: data.MRs{
				{
					Name:       "MR 1",
					MergedDate: helpers.Time(1047483647),
					Labels:     []string{"bugfix"},
				},
				{
					Name:       "MR 2",
					MergedDate: helpers.Time(1247483647),
				},
				{
					Name:       "MR 3",
					MergedDate: helpers.Time(1347483647),
					Labels:     []string{"enhancement"},
				},
			},
		},
		{
			name: "MRs with some matching exclude labels",
			args: args{
				m: data.MRs{
					{
						Name:       "MR 1",
						MergedDate: helpers.Time(1047483647),
						Labels:     []string{"bugfix", "enhancement"},
					},
					{
						Name:       "MR 2",
						MergedDate: helpers.Time(1247483647),
						Labels:     []string{"no changelog"},
					},
					{
						Name:       "MR 3",
						MergedDate: helpers.Time(1347483647),
						Labels:     []string{"enhancement"},
					},
					{
						Name:       "MR 4",
						MergedDate: helpers.Time(1347483647),
						Labels:     []string{"enhancement", "no changelog"},
					},
					{
						Name:       "MR 5",
						MergedDate: helpers.Time(1347483647),
						Labels:     []string{"bug", "wontfix"},
					},
				},
				withoutLabels: []string{"no changelog", "wontfix"},
			},
			want: data.MRs{
				{
					Name:       "MR 1",
					MergedDate: helpers.Time(1047483647),
					Labels:     []string{"bugfix", "enhancement"},
				},
				{
					Name:       "MR 3",
					MergedDate: helpers.Time(1347483647),
					Labels:     []string{"enhancement"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := data.FilterMRsByLabel(tt.args.m,
				tt.args.withoutLabels); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterMRsByLabel() = %v, want %v", got, tt.want)
			}
		})
	}
}
