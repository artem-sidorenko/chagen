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

import "time"

// Commit describes a struct with Commit information
type Commit struct {
	SHA          string
	AuthoredDate time.Time
	Title        string
}

// Commits returns different Commits
func Commits() []Commit {
	return []Commit{
		{"041152be02b2d69141d3a8d2278460f4777474f7", time.Unix(1047094647, 0),
			"Merge branch 'pr1' into 'master'"},
		{"1080a10971e4a887ae8a827bb16e0b04801f630b", time.Unix(1047194647, 0),
			"Merge branch 'pr2' into 'master'"},
		{"d72866aa0a25e58b7fb0365fba0fd6791d627451", time.Unix(1047294647, 0),
			"Merge branch 'pr3' into 'master'"},
		{"433a7f849f0a5c21a0f24886ff72a91e1e74888e", time.Unix(1047494647, 0),
			"Merge branch 'pr5' into 'master'"},
		{"e5bc67e0c5d2ed17639a6499d1d0c05d4073dc80", time.Unix(1047594647, 0),
			"Merge branch 'pr6' into 'master'"},
		{"d4c421f840e35fb15ae99683df23caf451db7377", time.Unix(1047694647, 0),
			"Merge branch 'pr7' into 'master'"},
		{"fd81ac08493e550604dd04fa39b9c2eb1907cea6", time.Unix(1047794647, 0),
			"Merge branch 'pr8' into 'master'"},
		{"cc1cf9b1441962bdd6b98a4e09363dffb2037835", time.Unix(1047894647, 0),
			"Merge branch 'pr9' into 'master'"},
		{"9772a06643b77ec1a16646df4bb909c771c09fba", time.Unix(1047994647, 0),
			"Merge branch 'pr10' into 'master'"},
		{"627b94d1e87e938ea140c592f3ebd115d5a98929", time.Unix(1048094647, 0),
			"Merge branch 'pr11' into 'master'"},
		{"c31af03759e2262d99b2c4a7571a8e0115f37d68", time.Unix(1048294647, 0),
			"Merge branch 'pr13' into 'master'"},
		{"9618c791ab1f643aeffb7c5e1abe5877223aaa91", time.Unix(1048394647, 0),
			"Merge branch 'pr14' into 'master'"},

		{"7d84cdb2f7c2d4619cda4b8adeb1897097b5c8fc", time.Unix(1047083647, 0), "Release v0.0.1"},
		{"b3622b516b8ad70ce5dc3fa422fb90c3b58fa9da", time.Unix(1047183647, 0), "Release v0.0.2"},
		{"52f214dc3bf6c0e2a87eae6eab363a317c5a665f", time.Unix(1047283647, 0), "Release v0.0.3"},
		{"d4ff341587bc80a9c897c28340df9fe8f9fc6309", time.Unix(1047383647, 0), "Release v0.0.4"},
		{"746e45ea014e257bcb7caa2c100ed1e5f63ed234", time.Unix(1047483647, 0), "Release v0.0.5"},
		{"ddde800c451bae606713ae0f8418badcf31db120", time.Unix(1047583647, 0), "Release v0.0.6"},
		{"d21438494dd0722c1d13dc496ae1f60fb85084c1", time.Unix(1047683647, 0), "Release v0.0.7"},
		{"8d8d817a530bc1c3f792d9508c187b5769c434c5", time.Unix(1047783647, 0), "Release v0.0.8"},
		{"fc9f16ecc043e3fe422834cd127311d11d423668", time.Unix(1047883647, 0), "Release v0.0.9"},
		{"dbbf36ffaae700a2ce03ef849d6f944031f34b95", time.Unix(1047983647, 0), "Release v0.1.0"},
		{"fc5d68ff1cf691e09f6ead044813274953c9b843", time.Unix(1048083647, 0), "Release v0.1.1"},
		{"d8351413f688c96c2c5d6fe58ebf5ac17f545bc0", time.Unix(1048183647, 0), "Release v0.1.2"},
	}
}

// CommitsBySHA returns a map with commit SHA as a key
func CommitsBySHA() map[string]Commit {
	rcommits := map[string]Commit{}

	for _, v := range Commits() {
		rcommits[v.SHA] = v
	}

	return rcommits
}
