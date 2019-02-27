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

	"github.com/artem-sidorenko/chagen/source/connectors/helpers"
	"github.com/google/go-github/github"
)

func genCommitAuthor(commitDate time.Time) *github.CommitAuthor {
	return &github.CommitAuthor{
		Date: getTimePtr(commitDate),
	}
}

func genCommit(sha string, commitDate time.Time) *github.Commit {
	return &github.Commit{
		SHA:       helpers.StringPtr(sha),
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
		Name:   helpers.StringPtr(name),
		Commit: genCommit(commitSha, commitDate),
	}
}

// nolint: unparam
func genRepositoryRelease(tagName, htmlURL string) *github.RepositoryRelease {
	return &github.RepositoryRelease{
		TagName: helpers.StringPtr(tagName),
		HTMLURL: helpers.StringPtr(htmlURL),
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
	labels []string,
) *github.Issue {

	var lbs []github.Label

	for _, l := range labels {
		lbs = append(lbs, github.Label{
			Name: helpers.StringPtr(l),
		})
	}

	return &github.Issue{
		Number:   getIntPtr(number),
		Title:    helpers.StringPtr(title),
		ClosedAt: getTimePtr(closedAt),
		HTMLURL:  helpers.StringPtr(htmlURL),
		Labels:   lbs,
	}
}

// nolint: unparam
func genIssuePR(
	number int, title string, prLink string,
) *github.Issue {
	return &github.Issue{
		Number:           getIntPtr(number),
		Title:            helpers.StringPtr(title),
		PullRequestLinks: &github.PullRequestLinks{URL: helpers.StringPtr(prLink)},
	}
}

// nolint: unparam
func genPR(
	number int,
	title, htmlURL, userLogin, userHTMLURL string,
	mergedAt time.Time, labels []string,
) *github.PullRequest {

	var lbs []*github.Label

	for _, l := range labels {
		lbs = append(lbs, &github.Label{
			Name: helpers.StringPtr(l),
		})
	}

	pr := &github.PullRequest{
		Number:  getIntPtr(number),
		Title:   helpers.StringPtr(title),
		HTMLURL: helpers.StringPtr(htmlURL),
		User: &github.User{
			Login:   helpers.StringPtr(userLogin),
			HTMLURL: helpers.StringPtr(userHTMLURL),
		},
		Labels: lbs,
	}

	if (mergedAt != time.Time{}) {
		pr.MergedAt = getTimePtr(mergedAt)
	}

	return pr
}
