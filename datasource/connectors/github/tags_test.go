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
	"context"
	"errors"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/artem-sidorenko/chagen/data"
	"github.com/artem-sidorenko/chagen/datasource/connectors/github"
	"github.com/artem-sidorenko/chagen/datasource/connectors/github/internal/testclient"
	"github.com/artem-sidorenko/chagen/internal/testing/helpers"
)

func TestConnector_Tags(t *testing.T) {
	tests := []struct {
		name        string
		returnValue testclient.ReturnValueStr
		want        data.Tags
		wantErr     error
		wantMaxtags []int
	}{
		{
			name: "API returns proper data",
			want: data.Tags{
				{
					Name:   "v0.1.2",
					Commit: "d8351413f688c96c2c5d6fe58ebf5ac17f545bc0",
					Date:   time.Unix(2048183647, 0),
					URL:    "https://github.com/testowner/testrepo/releases/v0.1.2",
				},
				{
					Name:   "v0.1.1",
					Commit: "fc5d68ff1cf691e09f6ead044813274953c9b843",
					Date:   time.Unix(2048083647, 0),
					URL:    "https://github.com/testowner/testrepo/releases/v0.1.1",
				},
				{
					Name:   "v0.1.0",
					Commit: "dbbf36ffaae700a2ce03ef849d6f944031f34b95",
					Date:   time.Unix(2047983647, 0),
					URL:    "https://github.com/testowner/testrepo/releases/v0.1.0",
				},
				{
					Name:   "v0.0.9",
					Commit: "fc9f16ecc043e3fe422834cd127311d11d423668",
					Date:   time.Unix(2047883647, 0),
					URL:    "https://github.com/testowner/testrepo/tree/v0.0.9",
				},
				{
					Name:   "v0.0.8",
					Commit: "8d8d817a530bc1c3f792d9508c187b5769c434c5",
					Date:   time.Unix(2047783647, 0),
					URL:    "https://github.com/testowner/testrepo/tree/v0.0.8",
				},
				{
					Name:   "v0.0.7",
					Commit: "d21438494dd0722c1d13dc496ae1f60fb85084c1",
					Date:   time.Unix(2047683647, 0),
					URL:    "https://github.com/testowner/testrepo/releases/v0.0.7",
				},
				{
					Name:   "v0.0.6",
					Commit: "ddde800c451bae606713ae0f8418badcf31db120",
					Date:   time.Unix(2047583647, 0),
					URL:    "https://github.com/testowner/testrepo/tree/v0.0.6",
				},
				{
					Name:   "v0.0.5",
					Commit: "746e45ea014e257bcb7caa2c100ed1e5f63ed234",
					Date:   time.Unix(2047483647, 0),
					URL:    "https://github.com/testowner/testrepo/tree/v0.0.5",
				},
				{
					Name:   "v0.0.4",
					Commit: "d4ff341587bc80a9c897c28340df9fe8f9fc6309",
					Date:   time.Unix(2047383647, 0),
					URL:    "https://github.com/testowner/testrepo/tree/v0.0.4",
				},
				{
					Name:   "v0.0.3",
					Commit: "52f214dc3bf6c0e2a87eae6eab363a317c5a665f",
					Date:   time.Unix(2047283647, 0),
					URL:    "https://github.com/testowner/testrepo/releases/v0.0.3",
				},
				{
					Name:   "v0.0.2",
					Commit: "b3622b516b8ad70ce5dc3fa422fb90c3b58fa9da",
					Date:   time.Unix(2047183647, 0),
					URL:    "https://github.com/testowner/testrepo/tree/v0.0.2",
				},
				{
					Name:   "v0.0.1",
					Commit: "7d84cdb2f7c2d4619cda4b8adeb1897097b5c8fc",
					Date:   time.Unix(2047083647, 0),
					URL:    "https://github.com/testowner/testrepo/releases/v0.0.1",
				},
			},
			wantMaxtags: []int{12},
		},
		{
			name: "ListTags call fails",
			returnValue: testclient.ReturnValueStr{
				RepoServiceListTagsErr: true,
			},
			wantErr: errors.New("can't fetch the tags"),
		},
		{
			name: "GetCommit call fails",
			returnValue: testclient.ReturnValueStr{
				RepoServiceGetCommitsErr: true,
			},
			wantErr: errors.New("can't fetch the commit"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			github.TagsPerPage = 5
			c := setupTestConnector(tt.returnValue, false)
			cerr := make(chan error, 1)

			cgot, _, cmaxtags := c.Tags(context.Background(), cerr)
			gotmaxtags := helpers.GetChannelValuesInt(cmaxtags)

			var got data.Tags
			for t := range cgot {
				got = append(got, t)
			}
			// sort the tags to have the stable order
			sort.Sort(&got)

			// sleep and allow the possible error to be delivered to the channel
			time.Sleep(time.Millisecond * 200)
			var err error
			select {
			case err = <-cerr:
			default:
			}

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Connector.Tags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connector.Tags() = %+v, want %+v", got, tt.want)
			}

			if err == nil { // compare the processed tags only in non-error situation
				if !reflect.DeepEqual(gotmaxtags, tt.wantMaxtags) {
					t.Errorf("Connector.Tags() maxtags = %v, want %v", gotmaxtags, tt.wantMaxtags)
				}
			}
		})
	}
}
