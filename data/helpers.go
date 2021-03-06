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

package data

func sliceContains(slice []string, str string) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}
	return false
}

// UTCDate sets all time within data to the UTC
func UTCDate(ts Tags, is Issues, mrs MRs) {
	for i, t := range ts {
		ts[i].Date = t.Date.UTC()
	}

	for i, d := range is {
		is[i].ClosedDate = d.ClosedDate.UTC()
	}

	for i, mr := range mrs {
		mrs[i].MergedDate = mr.MergedDate.UTC()
	}
}
