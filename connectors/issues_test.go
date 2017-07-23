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
	"github.com/artem-sidorenko/chagen/connectors"
	"reflect"
	"testing"
	"time"
)

func TestIssues_Sort(t *testing.T) {
	tests := []struct {
		name string
		is   *connectors.Issues
		want *connectors.Issues
	}{
		{
			name: "Issues are already sorted",
			is: &connectors.Issues{
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
			want: &connectors.Issues{
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
		},
		{
			name: "Issues are not sorted",
			is: &connectors.Issues{
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
			want: &connectors.Issues{
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.is.Sort()

			if !reflect.DeepEqual(tt.is, tt.want) {
				t.Errorf("Issues.Sort(), Issues = %v, want %v", tt.is, tt.want)
			}
		})
	}
}