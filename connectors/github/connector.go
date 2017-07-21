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

	"os"

	"net/url"
	"path"

	"github.com/artem-sidorenko/chagen/connectors"
)

// AccessTokenEnvVar contains the name of environment variable
// which sets the authentication access token
const AccessTokenEnvVar = "CHAGEN_GITHUB_TOKEN"

// Connector implements the GitHub connector
type Connector struct {
	API        API
	Owner      string
	Repo       string
	ProjectURL string
}

// Init takes the initialization of connector, e.g. reading environment vars etc
func (c *Connector) Init() {
	c.API = NewAPIClient(os.Getenv(AccessTokenEnvVar))
	c.Owner = "artem-sidorenko"
	c.Repo = "chef-cups"
	c.ProjectURL = "https://github.com/artem-sidorenko/chef-cups"
}

// GetTags returns the git tags
func (c *Connector) GetTags() (connectors.Tags, error) {
	tags, err := c.API.ListTags(c.Owner, c.Repo)
	if err != nil {
		return nil, err
	}

	var ret connectors.Tags
	for _, tag := range tags {
		tagName := tag.GetName()
		commit, err := c.API.GetCommit(c.Owner, c.Repo, tag.Commit.GetSHA())
		if err != nil {
			return nil, err
		}

		release, err := c.API.GetReleaseByTag(c.Owner, c.Repo, tagName)
		if err != nil {
			return nil, err
		}

		// if GitHub release for this tag was found -> use it
		// generate otherwise a link to the git tag view in the file tree
		var tagURL string
		if release != nil {
			tagURL = release.GetHTMLURL()
		} else {
			u, err := url.Parse(c.ProjectURL)
			if err != nil {
				return nil, err
			}
			u.Path = path.Join(u.Path, "/tree/"+tagName)
			tagURL = u.String()
		}

		ret = append(ret, connectors.Tag{
			Name:   tagName,
			Commit: commit.Commit.GetSHA(),
			Date:   commit.Commit.Committer.GetDate(),
			URL:    tagURL,
		})
	}
	return ret, nil
}

// GetIssues returns the closed issues
func (c *Connector) GetIssues() (connectors.Issues, error) {
	issues, err := c.API.ListIssues(c.Owner, c.Repo)
	if err != nil {
		return nil, err
	}

	var ret connectors.Issues
	for _, issue := range issues {
		//ensure we have an issue and not PR
		if issue.PullRequestLinks.GetURL() != "" {
			continue
		}

		ret = append(ret, connectors.Issue{
			ID:         issue.GetNumber(),
			Name:       issue.GetTitle(),
			ClosedDate: issue.GetClosedAt(),
			URL:        issue.GetHTMLURL(),
		})
	}

	return ret, nil
}

//GetMRs returns the merged pull requests
func (c *Connector) GetMRs() (connectors.MRs, error) {
	prs, err := c.API.ListPRs(c.Owner, c.Repo)
	if err != nil {
		return nil, err
	}

	var ret connectors.MRs
	for _, pr := range prs {
		// we need only merged PRs, skip everything else
		if pr.GetMergedAt() == (time.Time{}) {
			continue
		}

		ret = append(ret, connectors.MR{
			ID:         pr.GetNumber(),
			Name:       pr.GetTitle(),
			MergedDate: pr.GetMergedAt(),
		})
	}

	return ret, nil
}

func init() {
	connectors.RegisterConnector("github", "GitHub", &Connector{})
}
