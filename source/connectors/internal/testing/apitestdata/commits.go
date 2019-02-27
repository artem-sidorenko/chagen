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

package apitestdata

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
		{"7d84cdb2f7c2d4619cda4b8adeb1897097b5c8fc", time.Unix(2047083647, 0), "Release v0.0.1"},
		{"b3622b516b8ad70ce5dc3fa422fb90c3b58fa9da", time.Unix(2047183647, 0), "Release v0.0.2"},
		{"52f214dc3bf6c0e2a87eae6eab363a317c5a665f", time.Unix(2047283647, 0), "Release v0.0.3"},
		{"d4ff341587bc80a9c897c28340df9fe8f9fc6309", time.Unix(2047383647, 0), "Release v0.0.4"},
		{"746e45ea014e257bcb7caa2c100ed1e5f63ed234", time.Unix(2047483647, 0), "Release v0.0.5"},
		{"ddde800c451bae606713ae0f8418badcf31db120", time.Unix(2047583647, 0), "Release v0.0.6"},
		{"d21438494dd0722c1d13dc496ae1f60fb85084c1", time.Unix(2047683647, 0), "Release v0.0.7"},
		{"8d8d817a530bc1c3f792d9508c187b5769c434c5", time.Unix(2047783647, 0), "Release v0.0.8"},
		{"fc9f16ecc043e3fe422834cd127311d11d423668", time.Unix(2047883647, 0), "Release v0.0.9"},
		{"dbbf36ffaae700a2ce03ef849d6f944031f34b95", time.Unix(2047983647, 0), "Release v0.1.0"},
		{"fc5d68ff1cf691e09f6ead044813274953c9b843", time.Unix(2048083647, 0), "Release v0.1.1"},
		{"d8351413f688c96c2c5d6fe58ebf5ac17f545bc0", time.Unix(2048183647, 0), "Release v0.1.2"},
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
