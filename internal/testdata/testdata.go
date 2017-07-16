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

// Package testdata contains different test data
package testdata

import "github.com/artem-sidorenko/chagen/data"

// GeneratedChangelog contains the expected generated changelog
const GeneratedChangelog = `Changelog
=========

## [v0.1.0](https://example.com/release/v0.1.0) (2017-04-11)

Closed issues
-------------
- Test issue of new release [\#10](https://example.com/issue/10)

## [v0.0.1](https://example.com/release/v0.0.1) (2017-04-10)

Closed issues
-------------
- Test issue [\#1](https://example.com/issue/1)
`

// GeneratorDataStructure contains the data structure for changelog generator
var GeneratorDataStructure = []data.Release{
	{
		Release:    "v0.1.0",
		ReleaseURL: "https://example.com/release/v0.1.0",
		Date:       "2017-04-11",
		Issues: []data.Issue{
			{
				Name: "Test issue of new release",
				ID:   "10",
				URL:  "https://example.com/issue/10",
			},
		},
	},
	{
		Release:    "v0.0.1",
		ReleaseURL: "https://example.com/release/v0.0.1",
		Date:       "2017-04-10",
		Issues: []data.Issue{
			{
				Name: "Test issue",
				ID:   "1",
				URL:  "https://example.com/issue/1",
			},
		},
	},
}