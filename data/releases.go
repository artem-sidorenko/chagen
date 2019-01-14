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
	"sort"
	"time"
)

const releaseDateFormat = "02.01.2006"

// Release desribes a release with it data
type Release struct {
	Release    string
	ReleaseURL string
	Date       string
	Issues     Issues
	MRs        MRs
}

// Releases is a slice with Release elements
type Releases []Release

// NewReleases builds the Releases structure
// using given data from connector
func NewReleases(tags Tags, issues Issues, mrs MRs) Releases {
	var ret Releases
	var lastReleaseDate time.Time

	// we should have a proper sorted data to avoid surprises
	sort.Sort(&tags)
	sort.Sort(&issues)
	sort.Sort(&mrs)

	// as our tags are sorted, lets iterate from newest to the oldest
	for i, tag := range tags {
		// use the date of next tag (its older) as last release date
		// use 0 as last release date if we have the oldest (last) tag
		if i < (len(tags) - 1) {
			lastReleaseDate = tags[i+1].Date
		} else {
			lastReleaseDate = time.Time{}
		}

		ret = append(ret, Release{
			Release:    tag.Name,
			ReleaseURL: tag.URL,
			Date:       tag.Date.Format(releaseDateFormat),
			Issues:     issues.Filter(lastReleaseDate, tag.Date),
			MRs:        mrs.Filter(lastReleaseDate, tag.Date),
		})
	}

	return ret
}
