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

// Package testdata provides some test data
package testdata

import (
	"time"

	"github.com/artem-sidorenko/chagen/data"
)

// Tags returns some tags
func Tags() data.Tags {
	return data.Tags{
		{
			Name:   "v0.0.2",
			Date:   time.Unix(1147483647, 0),
			Commit: "b6a735dcb420a82865abe8c194900e59f6af9dea",
			URL:    "https://test.example.com/tags/v0.0.2",
		},
		{
			Name:   "v0.0.1",
			Date:   time.Unix(1047483647, 0),
			Commit: "d85645cbe6288cce5e5d066f8c7864040266cce3",
			URL:    "https://test.example.com/tags/v0.0.1",
		},
		{
			Name:   "v0.0.3",
			Date:   time.Unix(1247483647, 0),
			Commit: "25362c337d524025bf98e978059bf9bcd2b56221",
			URL:    "https://test.example.com/tags/v0.0.3",
		},
	}
}

// Issues returns some issues
func Issues() data.Issues {
	return data.Issues{
		{
			ID:         2,
			Name:       "Issue 2",
			ClosedDate: time.Unix(1247483647, 0),
			URL:        "http://test.example.com/issues/2",
		},
		{
			ID:         1,
			Name:       "Issue 1",
			ClosedDate: time.Unix(1047483647, 0),
			URL:        "http://test.example.com/issues/1",
		},
		{
			ID:         3,
			Name:       "Issue 3",
			ClosedDate: time.Unix(1347483647, 0),
			URL:        "http://test.example.com/issues/3",
		},
		{
			ID:         4,
			Name:       "Issue 4",
			ClosedDate: time.Unix(1297483647, 0),
			URL:        "http://test.example.com/issues/4",
		},
		{
			ID:         5,
			Name:       "Issue 5",
			ClosedDate: time.Unix(1298483647, 0),
			URL:        "http://test.example.com/issues/5",
			Labels:     []string{"issue5"},
		},
		{
			ID:         6,
			Name:       "Issue 6",
			ClosedDate: time.Unix(1299483647, 0),
			URL:        "http://test.example.com/issues/6",
			Labels:     []string{"no changelog"},
		},
	}
}
