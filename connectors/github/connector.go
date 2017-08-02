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
	"time"

	"errors"
	"os"

	"net/url"
	"path"

	"fmt"

	"github.com/artem-sidorenko/chagen/connectors"
	"github.com/artem-sidorenko/chagen/data"
	"github.com/urfave/cli"
)

// AccessTokenEnvVar contains the name of environment variable
// which sets the authentication access token
const AccessTokenEnvVar = "CHAGEN_GITHUB_TOKEN"

// Connector implements the GitHub connector
type Connector struct {
	API                 API
	Owner               string
	Repo                string
	ProjectURL          string
	NewTagUseReleaseURL bool
}

// Init takes the initialization of connector, e.g. reading environment vars etc
func (c *Connector) Init(cli *cli.Context) error {
	c.Owner = cli.String("github-owner")
	if c.Owner == "" {
		return errors.New("Option --github-owner is required")
	}
	c.Repo = cli.String("github-repo")
	if c.Repo == "" {
		return errors.New("Option --github-repo is required")
	}
	c.NewTagUseReleaseURL = cli.Bool("github-release-url")

	c.API = NewAPIClient(os.Getenv(AccessTokenEnvVar))
	c.ProjectURL = fmt.Sprintf("https://github.com/%s/%s", c.Owner, c.Repo)
	return nil
}

// GetNewTagURL returns the URL for a new tag, which does not exist yet
func (c *Connector) GetNewTagURL(TagName string) (string, error) {
	return c.getTagURL(TagName, c.NewTagUseReleaseURL)
}

// getTagURL returns the URL for a given tag.
// If alwaysUseReleaseURL is true: URL is provided for release page,
// even if it does not exist yet
func (c *Connector) getTagURL(TagName string, alwaysUseReleaseURL bool) (string, error) {
	release, err := c.API.GetReleaseByTag(c.Owner, c.Repo, TagName)
	if err != nil {
		return "", err
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
			u.Path = path.Join(u.Path, "/releases/"+TagName)
		} else { // build tag url
			u.Path = path.Join(u.Path, "/tree/"+TagName)
		}
		tagURL = u.String()
	}
	return tagURL, nil
}

// GetTags returns the git tags
func (c *Connector) GetTags() (data.Tags, error) {
	tags, err := c.API.ListTags(c.Owner, c.Repo)
	if err != nil {
		return nil, err
	}

	var ret data.Tags
	for _, tag := range tags {
		tagName := tag.GetName()
		commit, err := c.API.GetCommit(c.Owner, c.Repo, tag.Commit.GetSHA())
		if err != nil {
			return nil, err
		}

		tagURL, err := c.getTagURL(tagName, false)
		if err != nil {
			return nil, err
		}

		ret = append(ret, data.Tag{
			Name:   tagName,
			Commit: commit.Commit.GetSHA(),
			Date:   commit.Commit.Committer.GetDate(),
			URL:    tagURL,
		})
	}
	return ret, nil
}

// GetIssues returns the closed issues
func (c *Connector) GetIssues() (data.Issues, error) {
	issues, err := c.API.ListIssues(c.Owner, c.Repo)
	if err != nil {
		return nil, err
	}

	var ret data.Issues
	for _, issue := range issues {
		//ensure we have an issue and not PR
		if issue.PullRequestLinks.GetURL() != "" {
			continue
		}

		ret = append(ret, data.Issue{
			ID:         issue.GetNumber(),
			Name:       issue.GetTitle(),
			ClosedDate: issue.GetClosedAt(),
			URL:        issue.GetHTMLURL(),
		})
	}

	return ret, nil
}

//GetMRs returns the merged pull requests
func (c *Connector) GetMRs() (data.MRs, error) {
	prs, err := c.API.ListPRs(c.Owner, c.Repo)
	if err != nil {
		return nil, err
	}

	var ret data.MRs
	for _, pr := range prs {
		// we need only merged PRs, skip everything else
		if pr.GetMergedAt() == (time.Time{}) {
			continue
		}

		ret = append(ret, data.MR{
			ID:         pr.GetNumber(),
			Name:       pr.GetTitle(),
			MergedDate: pr.GetMergedAt(),
			URL:        pr.GetHTMLURL(),
			Author:     pr.User.GetLogin(),
			AuthorURL:  pr.User.GetHTMLURL(),
		})
	}

	return ret, nil
}

func init() {
	connectors.RegisterConnector("github", "GitHub", &Connector{}, []cli.Flag{
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
	})
}
