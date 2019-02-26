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

// Package testclient implements the used interfaces of github API client library
// and simulates the API answers for our tests
package testclient

import (
	"context"
	"fmt"

	"github.com/artem-sidorenko/chagen/source/connectors/github/internal/client"
	"github.com/artem-sidorenko/chagen/source/connectors/internal/testing/apitestdata"

	"github.com/google/go-github/github"
)

// ReturnValueStr represents the possible error controlling of API calls for testing
// if a field is set to true - return error, otherwise not
type ReturnValueStr struct {
	RetRepoServiceListTagsErr    bool
	RetRepoServiceGetCommitsErr  bool
	RetIssueServiceListByRepoErr bool
	RetPullRequestsListErr       bool
	RetRepoServiceGetErr         bool
	RetRepoServiceGetRespCode    int
}

// ReturnValue controls the error return values of API calls
// for testclient instances created by New
var ReturnValue = ReturnValueStr{} // nolint: gochecknoglobals

// GitHubRepoService simulates the github.RepositoriesService
type GitHubRepoService struct {
	RetRepositoryTags     []*github.RepositoryTag
	RetRepositoryCommits  map[string]*github.RepositoryCommit
	RetRepositoryReleases map[string]*github.RepositoryRelease
	ReturnValue           ReturnValueStr
}

// ListTags simulates the (github.RepositoriesService) ListTags call
func (g *GitHubRepoService) ListTags(
	ctx context.Context,
	owner, repo string,
	opt *github.ListOptions,
) ([]*github.RepositoryTag, *github.Response, error) {

	if g.ReturnValue.RetRepoServiceListTagsErr {
		return nil, nil, fmt.Errorf("can't fetch the tags")
	}

	resp, start, end := calcPaging(opt.Page, opt.PerPage, len(g.RetRepositoryTags))

	return g.RetRepositoryTags[start:end], resp, nil
}

// GetCommit simulates the (github.RepositoriesService) GetCommit call
func (g *GitHubRepoService) GetCommit(
	ctx context.Context,
	owner, repo, sha string,
) (*github.RepositoryCommit, *github.Response, error) {

	if g.ReturnValue.RetRepoServiceGetCommitsErr {
		return nil, nil, fmt.Errorf("can't fetch the commit")
	}

	if cm, ok := g.RetRepositoryCommits[sha]; ok {
		return cm, nil, nil
	}
	return nil, nil, fmt.Errorf("commit %v is not present", sha)
}

// GetReleaseByTag simulates the (github.RepositoriesService) GetCommit call
func (g *GitHubRepoService) GetReleaseByTag(
	ctx context.Context,
	owner, repo, tag string,
) (*github.RepositoryRelease, *github.Response, error) {
	if re, ok := g.RetRepositoryReleases[tag]; ok {
		return re, genResponse(200), nil
	}

	return nil, genResponse(404), nil
}

// Get simulates the (github.RepositoriesService) Get call
func (g *GitHubRepoService) Get(
	ctx context.Context,
	owner, repo string) (*github.Repository, *github.Response, error) {

	//if return code not defined, return 200 for Ok
	respCode := 200
	if g.ReturnValue.RetRepoServiceGetRespCode != 0 {
		respCode = g.ReturnValue.RetRepoServiceGetRespCode
	}

	response := genResponse(respCode)

	if g.ReturnValue.RetRepoServiceGetErr {
		return nil, response, fmt.Errorf("can't fetch the repo data")
	}

	return nil, response, nil
}

// GitHubIssueService simulates the github.IssuesService
type GitHubIssueService struct {
	RetIssues     []*github.Issue
	RetErrControl ReturnValueStr
}

// ListByRepo simulates the (github.IssuesService) ListByRepo call
func (g *GitHubIssueService) ListByRepo(
	ctx context.Context,
	owner string, repo string,
	opt *github.IssueListByRepoOptions,
) ([]*github.Issue, *github.Response, error) {

	if g.RetErrControl.RetIssueServiceListByRepoErr {
		return nil, nil, fmt.Errorf("can't fetch the issues")
	}

	resp, start, end := calcPaging(opt.Page, opt.PerPage, len(g.RetIssues))

	return g.RetIssues[start:end], resp, nil
}

// GitHubPullRequestsService simulates the github.PullRequestsService
type GitHubPullRequestsService struct {
	RetPRs        []*github.PullRequest
	RetErrControl ReturnValueStr
}

// List simulates the (github.PullRequestsService) ListByRepo call
func (g *GitHubPullRequestsService) List(
	ctx context.Context, owner string, repo string,
	opt *github.PullRequestListOptions,
) ([]*github.PullRequest, *github.Response, error) {

	if g.RetErrControl.RetPullRequestsListErr {
		return nil, nil, fmt.Errorf("can't fetch the PRs")
	}

	resp, start, end := calcPaging(opt.Page, opt.PerPage, len(g.RetPRs))

	return g.RetPRs[start:end], resp, nil
}

// newGitHubRepoService returns initialized instance of GitHubRepoService
func newGitHubRepoService() *GitHubRepoService {
	rtags := []*github.RepositoryTag{}
	rcommits := map[string]*github.RepositoryCommit{}
	rreleases := map[string]*github.RepositoryRelease{}

	for _, v := range apitestdata.Tags() {
		rtags = append(rtags, genRepositoryTag(v.Tag, v.Commit, v.Time))
		rcommits[v.Commit] = genRepositoryCommit(v.Commit, v.Time)
		if v.ReleasePresent {
			rreleases[v.Tag] = genRepositoryRelease(
				v.Tag,
				fmt.Sprintf("https://github.com/testowner/testrepo/releases/%v", v.Tag),
			)
		}
	}

	return &GitHubRepoService{
		ReturnValue:           ReturnValue,
		RetRepositoryTags:     rtags,
		RetRepositoryCommits:  rcommits,
		RetRepositoryReleases: rreleases,
	}
}

// newGitHubIssueService returns initialized instance of GitHubIssueService
func newGitHubIssueService() *GitHubIssueService {
	rissues := []*github.Issue{}

	for _, v := range apitestdata.Issues() {
		var i *github.Issue
		if v.PR {
			i = genIssuePR(
				v.ID, v.Title, fmt.Sprintf("https://example.com/prs/%v", v.ID),
			)
		} else {
			i = genIssue(
				v.ID, v.Title, v.ClosedAt,
				fmt.Sprintf("http://example.com/issues/%v", v.ID),
				v.Labels,
			)
		}
		rissues = append(rissues, i)
	}

	return &GitHubIssueService{
		RetErrControl: ReturnValue,
		RetIssues:     rissues,
	}
}

// newGitHubPullRequestsService returns initialized instance of GitHubPullRequestsService
// completely filled with provided testdata
func newGitHubPullRequestsService() *GitHubPullRequestsService {
	rprs := []*github.PullRequest{}

	for _, v := range apitestdata.MRs() {
		rprs = append(rprs, genPR(
			v.ID, v.Title,
			fmt.Sprintf("https://example.com/pulls/%v", v.ID),
			v.Username, fmt.Sprintf("https://example.com/users/%v", v.Username),
			v.MergedAt, v.Labels,
		))
	}

	return &GitHubPullRequestsService{
		RetErrControl: ReturnValue,
		RetPRs:        rprs,
	}
}

// New returns the configured simulated github API client
func New(_ context.Context, _ string) *client.GitHubClient {
	return &client.GitHubClient{
		Repositories: newGitHubRepoService(),
		Issues:       newGitHubIssueService(),
		PullRequests: newGitHubPullRequestsService(),
	}
}
