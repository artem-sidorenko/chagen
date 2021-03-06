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
	"time"
)

// MR describes a Pull or Merge Request
type MR struct {
	ID         int
	Name       string
	URL        string
	Author     string
	AuthorURL  string
	MergedDate time.Time
	Labels     []string
}

// MRs is a slice with MR elements
type MRs []MR

// Len implements the Sort.Interface
func (m *MRs) Len() int {
	return len(*m)
}

// Less implements the Sort.Interface
func (m *MRs) Less(i, j int) bool {
	return (*m)[i].MergedDate.After((*m)[j].MergedDate)
}

// Swap implements the Sort.Interface
func (m *MRs) Swap(i, j int) {
	(*m)[i], (*m)[j] = (*m)[j], (*m)[i]
}

// FilterMRs filters and returns new slice of Issues, where ClosedDate is between given dates
func FilterMRs(m MRs, fromDate, toDate time.Time) MRs {
	var ret MRs
	for _, mr := range m {
		if mr.MergedDate.After(fromDate) &&
			(mr.MergedDate.Before(toDate) || mr.MergedDate.Equal(toDate)) {
			ret = append(ret, mr)
		}
	}
	return ret
}

// FilterMRsByLabel filters out the MRs with given labels
func FilterMRsByLabel(m MRs, withoutLabels []string) MRs {
	var ret MRs
	for _, mr := range m {
		found := false
		for _, label := range mr.Labels {
			if sliceContains(withoutLabels, label) {
				found = true
				break
			}
		}

		if !found {
			ret = append(ret, mr)
		}
	}
	return ret
}
