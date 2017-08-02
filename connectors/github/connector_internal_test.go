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

package github

import (
	"reflect"
	"testing"

	"github.com/artem-sidorenko/chagen/connectors/github/internal/testapiclient"
	"github.com/google/go-github/github"
)

func TestConnector_getTagURL(t *testing.T) {
	type args struct {
		TagName             string
		alwaysUseReleaseURL bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "Release is present, alwaysUseReleaseURL is disabled",
			args: args{
				alwaysUseReleaseURL: false,
				TagName:             "v0.0.1",
			},
			want: "https://example.com/testowner/restrepo/releases/v0.0.1",
		},
		{
			name: "Release is present, alwaysUseReleaseURL is enabled",
			args: args{
				alwaysUseReleaseURL: false,
				TagName:             "v0.0.1",
			},
			want: "https://example.com/testowner/restrepo/releases/v0.0.1",
		},
		{
			name: "Release is not present, alwaysUseReleaseURL is disabled",
			args: args{
				alwaysUseReleaseURL: false,
				TagName:             "v0.0.3",
			},
			want: "https://example.com/testowner/restrepo/tree/v0.0.3",
		},
		{
			name: "Release is not present, alwaysUseReleaseURL is enabled",
			args: args{
				alwaysUseReleaseURL: true,
				TagName:             "v0.0.3",
			},
			want: "https://example.com/testowner/restrepo/releases/v0.0.3",
		},
	}
	for _, tt := range tests {
		c := &Connector{
			API: &testapiclient.TestAPIClient{
				RetGetReleaseByTag: map[string]*github.RepositoryRelease{
					"v0.0.1": {
						TagName: testapiclient.GetStringPtr("v0.0.1"),
						HTMLURL: testapiclient.GetStringPtr("https://example.com/testowner/restrepo/releases/v0.0.1"),
					},
				},
			},
			ProjectURL: "https://example.com/testowner/restrepo",
		}
		t.Run(tt.name, func(t *testing.T) {

			got, err := c.getTagURL(tt.args.TagName, tt.args.alwaysUseReleaseURL)
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("Connector.getTagURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connector.getTagURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
