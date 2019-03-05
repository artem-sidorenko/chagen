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

// Issue describes a struct with issue information
type Issue struct {
	ID       int
	Title    string
	ClosedAt time.Time
	Labels   []string
	PR       bool
}

// Issues returns different issues
func Issues() []Issue {
	return []Issue{
		{1214, "Test issue title 1", time.Unix(2047093647, 0), []string{"enhancement"}, false},
		{1227, "Test issue title 2", time.Unix(2047193647, 0), []string{"enhancement", "bugfix"}, false},
		{1239, "Test PR title 3", time.Unix(2047293647, 0), nil, true},
		{1244, "Test issue title 4", time.Unix(2047393647, 0), nil, false},
		{1254, "Test PR title 5", time.Unix(2047493647, 0), []string{"wontfix"}, true},
		{1264, "Test issue title 6", time.Unix(2047593647, 0), []string{"invalid"}, false},
		{1274, "Test issue title 7", time.Unix(2047693647, 0), []string{"no changelog"}, false},
		{1284, "Test PR title 8", time.Unix(2047793647, 0), []string{"enhancement"}, true},
		{1294, "Test issue title 9", time.Unix(2047893647, 0), nil, false},
		{1304, "Test issue title 10", time.Unix(2047993647, 0), []string{"wontfix"}, false},
		{1214, "Test PR title 11", time.Unix(2048093647, 0), []string{"enhancement"}, true},
		{1224, "Test issue title 12", time.Unix(2048193647, 0), nil, false},
		{1234, "Test issue title 13", time.Unix(2048293647, 0), []string{"enhancement"}, false},
	}
}
