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

	lastpage := len(g.RetRepositoryTags) / opt.PerPage
	if (len(g.RetRepositoryTags) % opt.PerPage) != 0 {
		lastpage++
	}

	nextpage := 0
	if opt.Page < lastpage {
		nextpage = opt.Page + 1
	}

	resp := &github.Response{
		NextPage: nextpage,
		LastPage: lastpage,
	}

	start := opt.PerPage * (opt.Page - 1)
	end := opt.PerPage * opt.Page
	if end > len(g.RetRepositoryTags) {
		end = len(g.RetRepositoryTags)
	}

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

	resp := &github.Response{
		NextPage: 0,
	}

	return g.RetIssues, resp, nil
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

	resp := &github.Response{
		NextPage: 0,
	}

	return g.RetPRs, resp, nil
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

	r := &GitHubRepoService{
		ReturnValue:           ReturnValue,
		RetRepositoryTags:     rtags,
		RetRepositoryCommits:  rcommits,
		RetRepositoryReleases: rreleases,
	}

	return r
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

	i := &GitHubIssueService{
		RetErrControl: ReturnValue,
		RetIssues: []*github.Issue{
			genIssue(
				1234, "Test issue title",
				time.Unix(1047483647, 0), "http://example.com/issues/1234",
				[]string{"enhancement"},
			),
			genIssuePR(4321, "Test PR title", "https://example.com/prs/4321"),
		},
	}

	p := &GitHubPullRequestsService{
		RetErrControl: ReturnValue,
		RetPRs: []*github.PullRequest{
			genPR(1234, "Test PR title", "https://example.com/pulls/1234",
				"test-user", "https://example.com/users/test-user",
				time.Unix(1747483647, 0), []string{"bugfix"}),
			genPR(1233, "Second closed PR title", "https://example.com/pulls/1233",
				"test-user", "https://example.com/users/test-user", time.Time{}, []string{}),
		},
	}

	return &client.GitHubClient{
		Repositories: r,
		Issues:       i,
		PullRequests: p,
	}
}
