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

package github_test

import (
	"testing"

	"github.com/artem-sidorenko/chagen/source/connectors/github/internal/testclient"
)

func TestGetNewTagURL(t *testing.T) {
	type fields struct {
		NewTagUseReleaseURL bool
	}
	type args struct {
		TagName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Release is present, NewTagUseReleaseURL is disabled",
			fields: fields{
				NewTagUseReleaseURL: false,
			},
			args: args{
				TagName: "v0.0.1",
			},
			want: "https://github.com/testowner/testrepo/releases/v0.0.1",
		},
		{
			name: "Release is present, NewTagUseReleaseURL is enabled",
			fields: fields{
				NewTagUseReleaseURL: true,
			},
			args: args{
				TagName: "v0.0.1",
			},
			want: "https://github.com/testowner/testrepo/releases/v0.0.1",
		},
		{
			name: "Release is not present, NewTagUseReleaseURL is disabled",
			fields: fields{
				NewTagUseReleaseURL: false,
			},
			args: args{
				TagName: "v0.2.3",
			},
			want: "https://github.com/testowner/testrepo/tree/v0.2.3",
		},
		{
			name: "Release is not present, alwaysUseReleaseURL is enabled",
			fields: fields{
				NewTagUseReleaseURL: true,
			},
			args: args{
				TagName: "v0.2.3",
			},
			want: "https://github.com/testowner/testrepo/releases/v0.2.3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestConnector(testclient.ReturnValueStr{}, tt.fields.NewTagUseReleaseURL)

			got, err := c.GetNewTagURL(tt.args.TagName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connector.GetNewTagURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Connector.GetNewTagURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
