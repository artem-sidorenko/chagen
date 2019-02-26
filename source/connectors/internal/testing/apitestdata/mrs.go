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

// MR describes a struct with PR/MR information
type MR struct {
	ID       int
	Title    string
	Username string
	MergedAt time.Time
	Labels   []string
}

// MRs returns different PRs/MRs
func MRs() []MR {
	return []MR{
		{2214, "Test PR title 1", "test-user", time.Unix(2047094647, 0), []string{"bugfix"}},
		{2224, "Test PR title 2", "test-user2", time.Unix(2047194647, 0), nil},
		{2234, "Test PR title 3", "test-user", time.Unix(2047294647, 0),
			[]string{"enhancement", "bugfix"}},
		{2244, "Test PR title 4 closed", "test-user", time.Time{}, []string{"wontfix"}},
		{2254, "Test PR title 5", "test-user", time.Unix(2047494647, 0), []string{"bugfix"}},
		{2264, "Test PR title 6", "test-user", time.Unix(2047594647, 0), []string{"enhancement"}},
		{2274, "Test PR title 7", "test5-user", time.Unix(2047694647, 0), []string{"bugfix"}},
		{2284, "Test PR title 8", "test-user", time.Unix(2047794647, 0), []string{"invalid"}},
		{2294, "Test PR title 9", "test-user", time.Unix(2047894647, 0), []string{"bugfix"}},
		{2304, "Test PR title 10", "test-user", time.Unix(2047994647, 0), []string{"bugfix"}},
		{2314, "Test PR title 11", "test-user8", time.Unix(2048094647, 0), []string{"no changelog"}},
		{2324, "Test PR title 12 closed", "test-user", time.Time{}, []string{"bugfix"}},
		{2334, "Test PR title 13", "test-user", time.Unix(2048294647, 0), []string{"bugfix"}},
		{2344, "Test PR title 14", "te77st-user", time.Unix(2048394647, 0), []string{"bugfix"}},
	}
}
