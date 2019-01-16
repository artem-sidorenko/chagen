/*
   Copyright 2019 Artem Sidorenko <artem@posteo.de>

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

package data

import "testing"

func Test_sliceContains(t *testing.T) {
	type args struct {
		slice []string
		str   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Slice contains the given string",
			args: args{
				slice: []string{"one", "two", "three"},
				str:   "one",
			},
			want: true,
		},
		{
			name: "Slice does not contain the given string",
			args: args{
				slice: []string{"one", "two", "three"},
				str:   "four",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sliceContains(tt.args.slice, tt.args.str); got != tt.want {
				t.Errorf("sliceContains() = %v, want %v", got, tt.want)
			}
		})
	}
}
