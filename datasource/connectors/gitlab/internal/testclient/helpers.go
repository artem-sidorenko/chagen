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

package testclient

import (
	gitlab "github.com/xanzy/go-gitlab"
)

// calcPaging calculates the paging options for simulation
func calcPaging(page, perPage, lenElements int) (resp *gitlab.Response, start, end int) {
	if perPage == 0 { // return all elements if no paging is requested
		return &gitlab.Response{
			NextPage:   0,
			TotalPages: 1,
		}, 0, lenElements
	}
	lastPage := lenElements / perPage
	// some elements are over full pages, we will have another non-complete page
	if (lenElements % perPage) != 0 {
		lastPage++
	}

	nextPage := 0
	if page < lastPage {
		nextPage = page + 1
	}

	resp = &gitlab.Response{
		NextPage:   nextPage,
		TotalPages: lastPage,
	}

	start = perPage * (page - 1)
	end = perPage * page
	if end > lenElements {
		end = lenElements
	}

	return resp, start, end
}
