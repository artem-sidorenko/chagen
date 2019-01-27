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
	"time"

	"github.com/artem-sidorenko/chagen/connectors/github/internal/client"
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

type gitHubRepoServiceInput struct {
	tag            string
	commit         string
	time           time.Time
	releasePresent bool
}

// newGitHubRepoService returns initialized instance of GitHubRepoService
// completely filled with provided testdata
func newGitHubRepoService(rsinput []gitHubRepoServiceInput) *GitHubRepoService {
	rtags := []*github.RepositoryTag{}
	rcommits := map[string]*github.RepositoryCommit{}
	rreleases := map[string]*github.RepositoryRelease{}
	for _, v := range rsinput {
		rtags = append(rtags, genRepositoryTag(v.tag, v.commit, v.time))
		rcommits[v.commit] = genRepositoryCommit(v.commit, v.time)
		if v.releasePresent {
			rreleases[v.tag] = genRepositoryRelease(
				v.tag,
				fmt.Sprintf("https://github.com/testowner/testrepo/releases/%v", v.tag),
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

type gitHubIssueServiceInput struct {
	id       int
	title    string
	closedAt time.Time
	labels   []string
	PR       bool
}

// newGitHubIssueService returns initialized instance of GitHubIssueService
// completely filled with provided testdata
func newGitHubIssueService(isinput []gitHubIssueServiceInput) *GitHubIssueService {
	rissues := []*github.Issue{}

	for _, v := range isinput {
		var i *github.Issue
		if v.PR {
			i = genIssuePR(
				v.id, v.title, fmt.Sprintf("https://example.com/prs/%v", v.id),
			)
		} else {
			i = genIssue(
				v.id, v.title, v.closedAt,
				fmt.Sprintf("http://example.com/issues/%v", v.id),
				v.labels,
			)
		}
		rissues = append(rissues, i)
	}

	return &GitHubIssueService{
		RetErrControl: ReturnValue,
		RetIssues:     rissues,
	}
}

type gitHubPullRequestsServiceInput struct {
	id       int
	title    string
	username string
	mergedAt time.Time
	labels   []string
}

// newGitHubPullRequestsService returns initialized instance of GitHubPullRequestsService
// completely filled with provided testdata
func newGitHubPullRequestsService(
	psinput []gitHubPullRequestsServiceInput,
) *GitHubPullRequestsService {

	rprs := []*github.PullRequest{}

	for _, v := range psinput {
		rprs = append(rprs, genPR(
			v.id, v.title,
			fmt.Sprintf("https://example.com/pulls/%v", v.id),
			v.username, fmt.Sprintf("https://example.com/users/%v", v.username),
			v.mergedAt, v.labels,
		))
	}

	return &GitHubPullRequestsService{
		RetErrControl: ReturnValue,
		RetPRs:        rprs,
	}
}

// New returns the configured simulated github API client
func New(_ context.Context, _ string) *client.GitHubClient {
	r := newGitHubRepoService([]gitHubRepoServiceInput{
		{"v0.0.1", "7d84cdb2f7c2d4619cda4b8adeb1897097b5c8fc", time.Unix(2047083647, 0), true},
		{"v0.0.2", "b3622b516b8ad70ce5dc3fa422fb90c3b58fa9da", time.Unix(2047183647, 0), false},
		{"v0.0.3", "52f214dc3bf6c0e2a87eae6eab363a317c5a665f", time.Unix(2047283647, 0), true},
		{"v0.0.4", "d4ff341587bc80a9c897c28340df9fe8f9fc6309", time.Unix(2047383647, 0), false},
		{"v0.0.5", "746e45ea014e257bcb7caa2c100ed1e5f63ed234", time.Unix(2047483647, 0), false},
		{"v0.0.6", "ddde800c451bae606713ae0f8418badcf31db120", time.Unix(2047583647, 0), false},
		{"v0.0.7", "d21438494dd0722c1d13dc496ae1f60fb85084c1", time.Unix(2047683647, 0), true},
		{"v0.0.8", "8d8d817a530bc1c3f792d9508c187b5769c434c5", time.Unix(2047783647, 0), false},
		{"v0.0.9", "fc9f16ecc043e3fe422834cd127311d11d423668", time.Unix(2047883647, 0), false},
		{"v0.1.0", "dbbf36ffaae700a2ce03ef849d6f944031f34b95", time.Unix(2047983647, 0), true},
		{"v0.1.1", "fc5d68ff1cf691e09f6ead044813274953c9b843", time.Unix(2048083647, 0), true},
		{"v0.1.2", "d8351413f688c96c2c5d6fe58ebf5ac17f545bc0", time.Unix(2048183647, 0), true},
	})

	i := newGitHubIssueService([]gitHubIssueServiceInput{
		{1214, "Test issue title 1", time.Unix(2047093647, 0), []string{"enhancement"}, false},
		{1227, "Test issue title 2", time.Unix(2047193647, 0), []string{"enhancement", "bugfix"}, false},
		{1239, "Test PR title 3", time.Unix(2047293647, 0), nil, true},
		{1244, "Test issue title 4", time.Unix(2047393647, 0), nil, false},
		{1254, "Test PR title 5", time.Unix(2047493647, 0), []string{"wontfix"}, true},
		{1264, "Test issue title 6", time.Unix(2047593647, 0), []string{"invalid"}, false},
		{1274, "Test issue title 7", time.Unix(2047693647, 0), []string{"no changelog"}, false},
		{1284, "Test PR title 8", time.Unix(2047793647, 0), []string{"enhancement"}, true},
		{1294, "Test issue title 9", time.Unix(2047893647, 0), []string{}, false},
		{1304, "Test issue title 10", time.Unix(2047993647, 0), []string{"wontfix"}, false},
		{1214, "Test PR title 11", time.Unix(2048093647, 0), []string{"enhancement"}, true},
		{1224, "Test issue title 12", time.Unix(2048193647, 0), nil, false},
		{1234, "Test issue title 13", time.Unix(2048293647, 0), []string{"enhancement"}, false},
	})

	p := newGitHubPullRequestsService([]gitHubPullRequestsServiceInput{
		{2214, "Test PR title 1", "test-user", time.Unix(2047094647, 0), []string{"bugfix"}},
		{2224, "Test PR title 2", "test-user2", time.Unix(2047194647, 0), nil},
		{2234, "Test PR title 3", "test-user", time.Unix(2047294647, 0),
			[]string{"enhancement", "bugfix"}},
		{2244, "Test PR title 4 closed", "test-user", time.Time{}, []string{"wontfix"}},
		{2254, "Test PR title 5", "test-user", time.Unix(2047494647, 0), []string{"bugfix"}},
		{2264, "Test PR title 6", "test-user", time.Unix(2047594647, 0), []string{"enhancement"}},
		{2274, "Test PR title 7", "test5-user", time.Unix(2047694647, 0), []string{"bugfix"}},
		{2284, "Test PR title 8", "test-user", time.Unix(2047794647, 0), []string{"invalid"}},
		{2294, "Test PR title 9", "test-user", time.Unix(2047894647, 0), []string{"bugfix"}},
		{2304, "Test PR title 10", "test-user", time.Unix(2047994647, 0), []string{"bugfix"}},
		{2314, "Test PR title 11", "test-user8", time.Unix(2048094647, 0), []string{"no changelog"}},
		{2324, "Test PR title 12 closed", "test-user", time.Time{}, []string{"bugfix"}},
		{2334, "Test PR title 13", "test-user", time.Unix(2048294647, 0), []string{"bugfix"}},
		{2344, "Test PR title 14", "te77st-user", time.Unix(2048394647, 0), []string{"bugfix"}},
	})

	return &client.GitHubClient{
		Repositories: r,
		Issues:       i,
		PullRequests: p,
	}
}
