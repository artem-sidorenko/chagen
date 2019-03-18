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
	"github.com/artem-sidorenko/chagen/internal/testing/helpers"
)

// MR describes a struct with PR/MR information
type MR struct {
	ID             int
	Title          string
	Username       string
	MergedAt       time.Time
	MergeCommitSHA string
	Labels         []string
}

// MRs returns different PRs/MRs
func MRs() []MR {
	return []MR{
		{2214, "Test PR title 1", "test-user", helpers.Time(1047094647),
			"041152be02b2d69141d3a8d2278460f4777474f7", []string{"bugfix"}},
		{2224, "Test PR title 2", "test-user2", helpers.Time(1047194647),
			"1080a10971e4a887ae8a827bb16e0b04801f630b", nil},
		{2234, "Test PR title 3", "test-user", helpers.Time(1047294647),
			"d72866aa0a25e58b7fb0365fba0fd6791d627451",
			[]string{"enhancement", "bugfix"}},
		{2244, "Test PR title 4 closed", "test-user", time.Time{},
			"", []string{"wontfix"}},
		{2254, "Test PR title 5", "test-user", helpers.Time(1047494647),
			"433a7f849f0a5c21a0f24886ff72a91e1e74888e", []string{"bugfix"}},
		{2264, "Test PR title 6", "test-user", helpers.Time(1047594647),
			"e5bc67e0c5d2ed17639a6499d1d0c05d4073dc80", []string{"enhancement"}},
		{2274, "Test PR title 7", "test5-user", helpers.Time(1047694647),
			"d4c421f840e35fb15ae99683df23caf451db7377", []string{"bugfix"}},
		{2284, "Test PR title 8", "test-user", helpers.Time(1047794647),
			"fd81ac08493e550604dd04fa39b9c2eb1907cea6", []string{"invalid"}},
		{2294, "Test PR title 9", "test-user", helpers.Time(1047894647),
			"cc1cf9b1441962bdd6b98a4e09363dffb2037835", []string{"bugfix"}},
		{2304, "Test PR title 10", "test-user", helpers.Time(1047994647),
			"9772a06643b77ec1a16646df4bb909c771c09fba", []string{"bugfix"}},
		{2314, "Test PR title 11", "test-user8", helpers.Time(1048094647),
			"627b94d1e87e938ea140c592f3ebd115d5a98929", []string{"no changelog"}},
		{2324, "Test PR title 12 closed", "test-user", time.Time{},
			"", []string{"bugfix"}},
		{2334, "Test PR title 13", "test-user", helpers.Time(1048294647),
			"c31af03759e2262d99b2c4a7571a8e0115f37d68", []string{"bugfix"}},
		{2344, "Test PR title 14", "te77st-user", helpers.Time(1048394647),
			"9618c791ab1f643aeffb7c5e1abe5877223aaa91", []string{"bugfix"}},
	}
}

// DataMRs returns the tags in the data.MR format
func DataMRs() []data.MR {
	var r []data.MR
	for _, m := range MRs() {
		r = append(r, data.MR{
			ID:         m.ID,
			Author:     m.Username,
			AuthorURL:  fmt.Sprintf("https://test.example.com/authors/%v", m.Username),
			Labels:     m.Labels,
			MergedDate: m.MergedAt,
			Name:       m.Title,
			URL:        fmt.Sprintf("https://test.example.com/mrs/%v", m.ID),
		})
	}
	return r
}
