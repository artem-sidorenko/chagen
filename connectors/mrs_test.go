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

package connectors_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/artem-sidorenko/chagen/connectors"
)

func TestMRs_Sort(t *testing.T) {
	tests := []struct {
		name string
		m    *connectors.MRs
		want *connectors.MRs
	}{
		{
			name: "MRs are already sorted",
			m: &connectors.MRs{
				{
					Name:       "MR 1",
					MergedDate: time.Unix(1047483647, 0),
				},
				{
					Name:       "MR 2",
					MergedDate: time.Unix(1247483647, 0),
				},
				{
					Name:       "MR 3",
					MergedDate: time.Unix(1347483647, 0),
				},
			},
			want: &connectors.MRs{
				{
					Name:       "MR 1",
					MergedDate: time.Unix(1047483647, 0),
				},
				{
					Name:       "MR 2",
					MergedDate: time.Unix(1247483647, 0),
				},
				{
					Name:       "MR 3",
					MergedDate: time.Unix(1347483647, 0),
				},
			},
		},
		{
			name: "MRs are not sorted",
			m: &connectors.MRs{
				{
					Name:       "MR 2",
					MergedDate: time.Unix(1247483647, 0),
				},
				{
					Name:       "MR 1",
					MergedDate: time.Unix(1047483647, 0),
				},
				{
					Name:       "MR 3",
					MergedDate: time.Unix(1347483647, 0),
				},
			},
			want: &connectors.MRs{
				{
					Name:       "MR 1",
					MergedDate: time.Unix(1047483647, 0),
				},
				{
					Name:       "MR 2",
					MergedDate: time.Unix(1247483647, 0),
				},
				{
					Name:       "MR 3",
					MergedDate: time.Unix(1347483647, 0),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Sort()

			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("MRs.Sort(), MRs = %v, want %v", tt.m, tt.want)
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
		m       *connectors.MRs
		args    args
		wantRet connectors.MRs
	}{
		{
			name: "Filtering of MRs",
			m: &connectors.MRs{
				{
					Name:       "MR 1",
					MergedDate: time.Unix(1047483647, 0),
				},
				{
					Name:       "MR 2",
					MergedDate: time.Unix(1247483647, 0),
				},
				{
					Name:       "MR 3",
					MergedDate: time.Unix(1347483647, 0),
				},
			},
			args: args{
				fromDate: time.Unix(1057483647, 0),
				toDate:   time.Unix(1337483647, 0),
			},
			wantRet: connectors.MRs{
				{
					Name:       "MR 2",
					MergedDate: time.Unix(1247483647, 0),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet := tt.m.Filter(tt.args.fromDate, tt.args.toDate)
			if !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("MRs.Filter() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}
