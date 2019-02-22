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
	"net/http"
	"time"

	gitlab "github.com/xanzy/go-gitlab"
)

func genTag(name string, commitSha string, commitDate time.Time) *gitlab.Tag {
	return &gitlab.Tag{
		Name:   name,
		Commit: genCommit(commitSha, commitDate),
	}
}

func genCommit(sha string, commitDate time.Time) *gitlab.Commit {
	return &gitlab.Commit{
		ID:           sha,
		AuthoredDate: &commitDate,
	}
}

func genResponse(statusCode int) *gitlab.Response {
	return &gitlab.Response{
		Response: &http.Response{StatusCode: statusCode},
	}
}
