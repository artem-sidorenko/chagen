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

	"github.com/google/go-github/github"
)

func genCommitAuthor(commitDate time.Time) *github.CommitAuthor {
	return &github.CommitAuthor{
		Date: getTimePtr(commitDate),
	}
}

func genCommit(sha string, commitDate time.Time) *github.Commit {
	return &github.Commit{
		SHA:       getStringPtr(sha),
		Committer: genCommitAuthor(commitDate),
	}
}

func genRepositoryCommit(sha string, commitDate time.Time) *github.RepositoryCommit {
	return &github.RepositoryCommit{
		Commit: genCommit(sha, commitDate),
	}
}

func genRepositoryTag(name string, commitSha string, commitDate time.Time) *github.RepositoryTag {
	return &github.RepositoryTag{
		Name:   getStringPtr(name),
		Commit: genCommit(commitSha, commitDate),
	}
}

// nolint: unparam
func genRepositoryRelease(tagName, htmlURL string) *github.RepositoryRelease {
	return &github.RepositoryRelease{
		TagName: getStringPtr(tagName),
		HTMLURL: getStringPtr(htmlURL),
	}
}

func genResponse(statusCode int) *github.Response {
	return &github.Response{
		Response: &http.Response{StatusCode: statusCode},
	}
}

// nolint: unparam
func genIssue(
	number int, title string,
	closedAt time.Time, htmlURL string,
) *github.Issue {

	return &github.Issue{
		Number:   getIntPtr(number),
		Title:    getStringPtr(title),
		ClosedAt: getTimePtr(closedAt),
		HTMLURL:  getStringPtr(htmlURL),
	}
}

// nolint: unparam
func genIssuePR(
	number int, title string, prLink string,
) *github.Issue {
	return &github.Issue{
		Number:           getIntPtr(number),
		Title:            getStringPtr(title),
		PullRequestLinks: &github.PullRequestLinks{URL: getStringPtr(prLink)},
	}
}

// nolint: unparam
func genPR(
	number int,
	title, htmlURL, userLogin, userHTMLURL string,
	mergedAt time.Time,
) *github.PullRequest {

	pr := &github.PullRequest{
		Number:  getIntPtr(number),
		Title:   getStringPtr(title),
		HTMLURL: getStringPtr(htmlURL),
		User: &github.User{
			Login:   getStringPtr(userLogin),
			HTMLURL: getStringPtr(userHTMLURL),
		},
	}

	if (mergedAt != time.Time{}) {
		pr.MergedAt = getTimePtr(mergedAt)
	}

	return pr
}
