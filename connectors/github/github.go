/*
   Copyright 2017 Artem Sidorenko <artem@posteo.de>

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

// Package github implements a github connector
package github

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/artem-sidorenko/chagen/connectors"
	"github.com/artem-sidorenko/chagen/connectors/github/internal/client"
	"github.com/artem-sidorenko/chagen/data"

	"github.com/google/go-github/github"
	"github.com/urfave/cli"
)

// AccessTokenEnvVar contains the name of environment variable
// which sets the authentication access token
const AccessTokenEnvVar = "CHAGEN_GITHUB_TOKEN" // nolint: gosec

// Connector implements the GitHub connector
type Connector struct {
	context             context.Context
	client              *client.GitHubClient
	Owner               string
	Repo                string
	ProjectURL          string
	NewTagUseReleaseURL bool
}

// NewGitHubClientFunc links to the constructor, which is used to create Connector.client
var NewGitHubClientFunc = client.NewGitHubClient // nolint: gochecknoglobals

// formatErrorCode formats the error message for this connector
func formatErrorCode(query string, err error) error {
	return fmt.Errorf("GitHub query '%s' failed: %s", query, err)
}

// RepositoryExists checks if referenced repository is present
func (c *Connector) RepositoryExists() (bool, error) {
	_, resp, err := c.client.Repositories.Get(c.context, c.Owner, c.Repo)
	if err != nil {
		if resp.StatusCode == 404 { // not found isn't an error
			return false, nil
		}
		return false, formatErrorCode("RepositoryExists", err)
	}
	switch resp.StatusCode {
	case 200:
		return true, nil
	default:
		return false, formatErrorCode(
			"RepositoryExists",
			fmt.Errorf("unhandled HTTP response code %v", resp.StatusCode),
		)
	}
}

// GetNewTagURL returns the URL for a new tag, which does not exist yet
func (c *Connector) GetNewTagURL(TagName string) (string, error) {
	return c.getTagURL(TagName, c.NewTagUseReleaseURL)
}

// getTagURL returns the URL for a given tag.
// If alwaysUseReleaseURL is true: URL is provided for release page,
// even if it does not exist yet
func (c *Connector) getTagURL(tagName string, alwaysUseReleaseURL bool) (string, error) {
	release, resp, err := c.client.Repositories.GetReleaseByTag(c.context, c.Owner, c.Repo, tagName)
	if err != nil {
		// no release was found for this tag, this is no error for us
		if resp.StatusCode != 404 {
			return "", formatErrorCode("getTagURL", err)
		}
	}

	// if GitHub release for this tag was found -> use it
	// generate otherwise a link to the git tag view in the file tree
	var tagURL string
	if release != nil { // we got real release URL, use it
		tagURL = release.GetHTMLURL()
	} else { // build own URL
		u, err := url.Parse(c.ProjectURL)
		if err != nil {
			return "", err
		}

		if alwaysUseReleaseURL { // try to build own release url
			u.Path = path.Join(u.Path, "/releases/"+tagName)
		} else { // build tag url
			u.Path = path.Join(u.Path, "/tree/"+tagName)
		}
		tagURL = u.String()
	}
	return tagURL, nil
}

// GetTags returns the git tags
func (c *Connector) GetTags() (data.Tags, error) {
	opt := &github.ListOptions{}

	var ret data.Tags
	for {
		tags, resp, err := c.client.Repositories.ListTags(c.context, c.Owner, c.Repo, nil)
		if err != nil {
			return nil, formatErrorCode("GetTags", err)
		}

		for _, tag := range tags {
			tagName := tag.GetName()

			commit, _, err := c.client.Repositories.GetCommit(c.context,
				c.Owner, c.Repo, tag.Commit.GetSHA())
			if err != nil {
				return nil, formatErrorCode("GetTags", err)
			}

			tagURL, err := c.getTagURL(tagName, false)
			if err != nil {
				return nil, formatErrorCode("GetTags", err)
			}

			ret = append(ret, data.Tag{
				Name:   tagName,
				Commit: commit.Commit.GetSHA(),
				Date:   commit.Commit.Committer.GetDate(),
				URL:    tagURL,
			})
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return ret, nil
}

// GetIssues returns the closed issues
func (c *Connector) GetIssues() (data.Issues, error) {
	opt := &github.IssueListByRepoOptions{State: "closed"}

	var ret data.Issues
	for {
		issues, resp, err := c.client.Issues.ListByRepo(c.context, c.Owner, c.Repo, opt)
		if err != nil {
			return nil, formatErrorCode("GetIssues", err)
		}

		for _, issue := range issues {
			//ensure we have an issue and not PR
			if issue.PullRequestLinks.GetURL() != "" {
				continue
			}

			var lbs []string
			if issue.Labels != nil && len(issue.Labels) > 0 {
				for _, l := range issue.Labels {
					lbs = append(lbs, *l.Name)
				}
			}

			ret = append(ret, data.Issue{
				ID:         issue.GetNumber(),
				Name:       issue.GetTitle(),
				ClosedDate: issue.GetClosedAt(),
				URL:        issue.GetHTMLURL(),
				Labels:     lbs,
			})
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return ret, nil
}

//GetMRs returns the merged pull requests
func (c *Connector) GetMRs() (data.MRs, error) {
	opt := &github.PullRequestListOptions{State: "closed"}

	var ret data.MRs
	for {
		prs, resp, err := c.client.PullRequests.List(c.context, c.Owner, c.Repo, opt)
		if err != nil {
			return nil, formatErrorCode("GetMRs", err)
		}

		for _, pr := range prs {
			// we need only merged PRs, skip everything else
			if pr.GetMergedAt() == (time.Time{}) {
				continue
			}

			var lbs []string
			if pr.Labels != nil && len(pr.Labels) > 0 {
				for _, l := range pr.Labels {
					lbs = append(lbs, *l.Name)
				}
			}

			ret = append(ret, data.MR{
				ID:         pr.GetNumber(),
				Name:       pr.GetTitle(),
				MergedDate: pr.GetMergedAt(),
				URL:        pr.GetHTMLURL(),
				Author:     pr.User.GetLogin(),
				AuthorURL:  pr.User.GetHTMLURL(),
				Labels:     lbs,
			})
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return ret, nil
}

// New returns a new initialized Connector or error if any
func New(ctx *cli.Context) (connectors.Connector, error) {
	owner := ctx.String("github-owner")
	if owner == "" {
		return nil, errors.New("option --github-owner is required")
	}
	repo := ctx.String("github-repo")
	if repo == "" {
		return nil, errors.New("option --github-repo is required")
	}
	newTagUseReleaseURL := ctx.Bool("github-release-url")

	return &Connector{
		context:             context.Background(),
		client:              NewGitHubClientFunc(context.Background(), os.Getenv(AccessTokenEnvVar)),
		Owner:               owner,
		Repo:                repo,
		NewTagUseReleaseURL: newTagUseReleaseURL,
		ProjectURL:          fmt.Sprintf("https://github.com/%s/%s", owner, repo),
	}, nil
}

// CLIFlags returns the possible CLI flags for this connector
func CLIFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "github-owner",
			Usage: "Owner/organisation where repository belongs to",
		},
		cli.StringFlag{
			Name:  "github-repo",
			Usage: "Name of repository",
		},
		cli.BoolFlag{
			Name:  "github-release-url",
			Usage: "New release should use URL to the GitHub release, even if it does not exist yet",
		},
	}
}

func init() { // nolint: gochecknoinits
	connectors.RegisterConnector("github", "GitHub", New, CLIFlags)
}
