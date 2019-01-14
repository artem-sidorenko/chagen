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

// Tag describes a git tag
type Tag struct {
	Name   string
	Commit string
	Date   time.Time
	URL    string
}

// Tags is a slice with Tag elements
type Tags []Tag

// Len implements the Sort.Interface
func (t *Tags) Len() int {
	return len(*t)
}

// Less implements the Sort.Interface
func (t *Tags) Less(i, j int) bool {
	return (*t)[i].Date.After((*t)[j].Date)
}

// Swap implements the Sort.Interface
func (t *Tags) Swap(i, j int) {
	(*t)[i], (*t)[j] = (*t)[j], (*t)[i]
}
