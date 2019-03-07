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

	"github.com/artem-sidorenko/chagen/datasource/connectors/github/internal/client"
	"github.com/artem-sidorenko/chagen/datasource/connectors/internal/testing/apitestdata"

	"github.com/google/go-github/github"
)

// ReturnValueStr represents the possible error controlling of API calls for testing
// if a field is set to true - return error, otherwise not
type ReturnValueStr struct {
	RepoServiceListTagsErr    bool
	RepoServiceGetCommitsErr  bool
	IssueServiceListByRepoErr bool
	PullRequestsListErr       bool
	RepoServiceGetErr         bool
	RepoServiceGetRespCode    int
}

// ReturnValue controls the error return values of API calls
// for testclient instances created by New
var ReturnValue = ReturnValueStr{} // nolint: gochecknoglobals

// RepoService simulates the github.RepositoriesService
type RepoService struct {
	RepositoryTags     []*github.RepositoryTag
	RepositoryCommits  map[string]*github.RepositoryCommit
	RepositoryReleases map[string]*github.RepositoryRelease
	ReturnValue        ReturnValueStr
}

// ListTags simulates the (github.RepositoriesService) ListTags call
func (g *RepoService) ListTags(
	ctx context.Context,
	owner, repo string,
	opt *github.ListOptions,
) ([]*github.RepositoryTag, *github.Response, error) {

	if g.ReturnValue.RepoServiceListTagsErr {
		return nil, nil, fmt.Errorf("can't fetch the tags")
	}

	resp, start, end := calcPaging(opt.Page, opt.PerPage, len(g.RepositoryTags))

	return g.RepositoryTags[start:end], resp, nil
}

// GetCommit simulates the (github.RepositoriesService) GetCommit call
func (g *RepoService) GetCommit(
	ctx context.Context,
	owner, repo, sha string,
) (*github.RepositoryCommit, *github.Response, error) {

	if g.ReturnValue.RepoServiceGetCommitsErr {
		return nil, nil, fmt.Errorf("can't fetch the commit")
	}

	if cm, ok := g.RepositoryCommits[sha]; ok {
		return cm, nil, nil
	}
	return nil, nil, fmt.Errorf("commit %v is not present", sha)
}

// GetReleaseByTag simulates the (github.RepositoriesService) GetCommit call
func (g *RepoService) GetReleaseByTag(
	ctx context.Context,
	owner, repo, tag string,
) (*github.RepositoryRelease, *github.Response, error) {
	if re, ok := g.RepositoryReleases[tag]; ok {
		return re, genResponse(200), nil
	}

	return nil, genResponse(404), nil
}

// Get simulates the (github.RepositoriesService) Get call
func (g *RepoService) Get(
	ctx context.Context,
	owner, repo string) (*github.Repository, *github.Response, error) {

	//if return code not defined, return 200 for Ok
	respCode := 200
	if g.ReturnValue.RepoServiceGetRespCode != 0 {
		respCode = g.ReturnValue.RepoServiceGetRespCode
	}

	response := genResponse(respCode)

	if g.ReturnValue.RepoServiceGetErr {
		return nil, response, fmt.Errorf("can't fetch the repo data")
	}

	return nil, response, nil
}

// IssueService simulates the github.IssuesService
type IssueService struct {
	Issues      []*github.Issue
	ReturnValue ReturnValueStr
}

// ListByRepo simulates the (github.IssuesService) ListByRepo call
func (g *IssueService) ListByRepo(
	ctx context.Context,
	owner string, repo string,
	opt *github.IssueListByRepoOptions,
) ([]*github.Issue, *github.Response, error) {

	if g.ReturnValue.IssueServiceListByRepoErr {
		return nil, nil, fmt.Errorf("can't fetch the issues")
	}

	resp, start, end := calcPaging(opt.Page, opt.PerPage, len(g.Issues))

	return g.Issues[start:end], resp, nil
}

// PullRequestsService simulates the github.PullRequestsService
type PullRequestsService struct {
	PRs         []*github.PullRequest
	ReturnValue ReturnValueStr
}

// List simulates the (github.PullRequestsService) ListByRepo call
func (g *PullRequestsService) List(
	ctx context.Context, owner string, repo string,
	opt *github.PullRequestListOptions,
) ([]*github.PullRequest, *github.Response, error) {

	if g.ReturnValue.PullRequestsListErr {
		return nil, nil, fmt.Errorf("can't fetch the PRs")
	}

	resp, start, end := calcPaging(opt.Page, opt.PerPage, len(g.PRs))

	return g.PRs[start:end], resp, nil
}

// newGitHubRepoService returns initialized instance of GitHubRepoService
func newGitHubRepoService() *RepoService {
	rtags := []*github.RepositoryTag{}
	rcommits := map[string]*github.RepositoryCommit{}
	rreleases := map[string]*github.RepositoryRelease{}

	for _, v := range apitestdata.Commits() {
		rcommits[v.SHA] = genRepositoryCommit(v.SHA, v.AuthoredDate)
	}

	for _, v := range apitestdata.Tags() {
		rtags = append(rtags, genRepositoryTag(v.Tag, rcommits[v.Commit].Commit))

		if v.ReleaseTime != nil {
			rreleases[v.Tag] = genRepositoryRelease(
				v.Tag,
				fmt.Sprintf("https://github.com/testowner/testrepo/releases/%v", v.Tag),
			)
		}
	}

	return &RepoService{
		ReturnValue:        ReturnValue,
		RepositoryTags:     rtags,
		RepositoryCommits:  rcommits,
		RepositoryReleases: rreleases,
	}
}

// newGitHubIssueService returns initialized instance of GitHubIssueService
func newGitHubIssueService() *IssueService {
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

	return &IssueService{
		ReturnValue: ReturnValue,
		Issues:      rissues,
	}
}

// newGitHubPullRequestsService returns initialized instance of GitHubPullRequestsService
// completely filled with provided testdata
func newGitHubPullRequestsService() *PullRequestsService {
	rprs := []*github.PullRequest{}

	for _, v := range apitestdata.MRs() {
		rprs = append(rprs, genPR(
			v.ID, v.Title,
			fmt.Sprintf("https://example.com/pulls/%v", v.ID),
			v.Username, fmt.Sprintf("https://example.com/users/%v", v.Username),
			v.MergedAt, v.Labels,
		))
	}

	return &PullRequestsService{
		ReturnValue: ReturnValue,
		PRs:         rprs,
	}
}

// New returns the configured simulated github API client
func New(_ context.Context, _ string) *client.Client {
	return &client.Client{
		Repositories: newGitHubRepoService(),
		Issues:       newGitHubIssueService(),
		PullRequests: newGitHubPullRequestsService(),
	}
}
