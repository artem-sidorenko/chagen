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

package connectors

import (
	"sort"
	"time"
)

// Issue describes an issue in the bug tracker
type Issue struct {
	ID         int
	Name       string
	ClosedDate time.Time
	URL        string
}

// Issues is a slice with Issue elements
type Issues []Issue

// Len implements the Sort.Interface
func (is *Issues) Len() int {
	return len(*is)
}

// Less implements the Sort.Interface
func (is *Issues) Less(i, j int) bool {
	return (*is)[i].ClosedDate.Before((*is)[j].ClosedDate)
}

// Swap implements the Sort.Interface
func (is *Issues) Swap(i, j int) {
	(*is)[i], (*is)[j] = (*is)[j], (*is)[i]
}

// Sort implements sorting of available Issues
func (is *Issues) Sort() {
	sort.Sort(is)
}

// Filter filters and returns new slice of Issues, where ClosedDate is between given dates
func (is *Issues) Filter(fromDate, toDate time.Time) (ret Issues) {
	for _, issue := range *is {
		if issue.ClosedDate.After(fromDate) && issue.ClosedDate.Before(toDate) {
			ret = append(ret, issue)
		}
	}
	return
}
