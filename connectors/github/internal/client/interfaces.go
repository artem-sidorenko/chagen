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

package client

import (
	"context"

	"github.com/google/go-github/github"
)

// inspired by https://github.com/google/go-github/issues/113
// lets have interfaces and structs which allow easy testing

// GithubRepoService describes the methods we use from
// github.RepositoriesService
type GithubRepoService interface {
	ListTags(
		ctx context.Context,
		owner, repo string,
		opt *github.ListOptions) ([]*github.RepositoryTag, *github.Response, error)
	GetCommit(
		ctx context.Context,
		owner, repo, sha string) (*github.RepositoryCommit, *github.Response, error)
	GetReleaseByTag(
		ctx context.Context,
		owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error)
}

// GithubIssuesService describes the methods we use from
// github.IssuesService
type GithubIssuesService interface {
	ListByRepo(
		ctx context.Context,
		owner string, repo string,
		opt *github.IssueListByRepoOptions) ([]*github.Issue, *github.Response, error)
}

// GitHubPullRequestsService describes the methods we use from
// github.PullRequestsService
type GitHubPullRequestsService interface {
	List(
		ctx context.Context, owner string, repo string,
		opt *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
}

// GitHubClient wraps the github.Client with interfaces we are using
type GitHubClient struct {
	Repositories GithubRepoService
	Issues       GithubIssuesService
	PullRequests GitHubPullRequestsService
}
