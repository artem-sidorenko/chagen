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

package data

import (
	"github.com/artem-sidorenko/chagen/connectors"
)

// Release desribes a release with it data
type Release struct {
	Release    string
	ReleaseURL string
	Date       string
	Issues     Issues
}

// Releases is a slice with Release elements
type Releases []Release

// NewReleases builds the Releases structure
// using given data from connector
func NewReleases(
	tags connectors.Tags,
	issues connectors.Issues) (ret Releases) {

	for _, tag := range tags {
		var relIssues Issues

		for _, is := range issues {
			relIssues = append(relIssues, Issue{
				ID:   is.ID,
				Name: is.Name,
			})
		}

		rel := Release{
			Release:    tag.Name,
			ReleaseURL: "",
			Date:       tag.Date.Format("02.01.2006"),
			Issues:     relIssues,
		}
		ret = append(ret, rel)
	}
	return
}
