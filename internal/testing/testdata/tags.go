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

package testdata

import (
	"fmt"
	"time"

	"github.com/artem-sidorenko/chagen/data"
	"github.com/artem-sidorenko/chagen/datasource/connectors/helpers"
)

// Tag describes a struct with tag information
type Tag struct {
	Tag         string
	Commit      string
	ReleaseTime *time.Time // if its nil -> there is no Release present, just the tag
}

// Tags returns different tags
func Tags() []Tag {
	return []Tag{
		{"v0.0.1", "7d84cdb2f7c2d4619cda4b8adeb1897097b5c8fc", helpers.TimePtr(time.Unix(2047083657, 0))},
		{"v0.0.2", "b3622b516b8ad70ce5dc3fa422fb90c3b58fa9da", nil},
		{"v0.0.3", "52f214dc3bf6c0e2a87eae6eab363a317c5a665f", helpers.TimePtr(time.Unix(2047283657, 0))},
		{"v0.0.4", "d4ff341587bc80a9c897c28340df9fe8f9fc6309", nil},
		{"v0.0.5", "746e45ea014e257bcb7caa2c100ed1e5f63ed234", nil},
		{"v0.0.6", "ddde800c451bae606713ae0f8418badcf31db120", nil},
		{"v0.0.7", "d21438494dd0722c1d13dc496ae1f60fb85084c1", helpers.TimePtr(time.Unix(2047683657, 0))},
		{"v0.0.8", "8d8d817a530bc1c3f792d9508c187b5769c434c5", nil},
		{"v0.0.9", "fc9f16ecc043e3fe422834cd127311d11d423668", nil},
		{"v0.1.0", "dbbf36ffaae700a2ce03ef849d6f944031f34b95", helpers.TimePtr(time.Unix(2047983657, 0))},
		{"v0.1.1", "fc5d68ff1cf691e09f6ead044813274953c9b843", helpers.TimePtr(time.Unix(2048083657, 0))},
		{"v0.1.2", "d8351413f688c96c2c5d6fe58ebf5ac17f545bc0", helpers.TimePtr(time.Unix(2048183657, 0))},
	}
}

// DataTags returns the tags in the data.Tag format
func DataTags() []data.Tag {
	var r []data.Tag
	commits := CommitsBySHA()
	for _, t := range Tags() {
		r = append(r, data.Tag{
			Commit: t.Commit,
			Date:   commits[t.Commit].AuthoredDate,
			Name:   t.Tag,
			URL:    fmt.Sprintf("https://test.example.com/tags/%v", t.Tag),
		})
	}
	return r
}
