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

package helpers

import "time"

// Time returns a UTC time for given timestamp
func Time(sec int64) time.Time {
	return time.Unix(sec, 0).UTC()
}

// TimePtr returns a pointer to time.Time for a given timestamp
func TimePtr(sec int64) *time.Time {
	t := Time(sec)
	return &t
}
