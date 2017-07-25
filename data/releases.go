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

// Release desribes a release with it data
type Release struct {
	Release       string
	ReleaseURL    string
	DateFormatted string
	Date          time.Time
	Issues        Issues
	MRs           MRs
}

// Releases is a slice with Release elements
type Releases []Release

// Len implements the Sort.Interface
func (r *Releases) Len() int {
	return len(*r)
}

// Less implements the Sort.Interface
func (r *Releases) Less(i, j int) bool {
	return (*r)[i].Date.Before((*r)[j].Date)
}

// Swap implements the Sort.Interface
func (r *Releases) Swap(i, j int) {
	(*r)[i], (*r)[j] = (*r)[j], (*r)[i]
}

// Sort implements reverse sorting by date (the oldest release is first)
func (r *Releases) Sort() {
	sort.Sort(sort.Reverse(r))
}

// NewReleases builds the Releases structure
// using given data from connector
func NewReleases(
	tags Tags,
	issues Issues,
	mrs MRs) (ret Releases) {

	var lastReleaseDate time.Time

	for _, tag := range tags {
		ret = append(ret, Release{
			Release:       tag.Name,
			ReleaseURL:    tag.URL,
			Date:          tag.Date,
			DateFormatted: tag.Date.Format("02.01.2006"),
			Issues:        issues.Filter(lastReleaseDate, tag.Date),
			MRs:           mrs.Filter(lastReleaseDate, tag.Date),
		})

		lastReleaseDate = tag.Date
	}

	return
}
