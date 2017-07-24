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

func TestTags_Sort(t *testing.T) {
	tests := []struct {
		name string
		t    *connectors.Tags
		want *connectors.Tags
	}{
		{
			name: "Tags are already sorted",
			t: &connectors.Tags{
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
			want: &connectors.Tags{
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
		},

		{
			name: "Tags are not sorted",
			t: &connectors.Tags{
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
			want: &connectors.Tags{
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.Sort()

			if !reflect.DeepEqual(tt.t, tt.want) {
				t.Errorf("Tags.Sort(), Tags = %v, want %v", tt.t, tt.want)
			}
		})
	}
}
