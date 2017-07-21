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

	"github.com/artem-sidorenko/chagen/connectors"
)

// Connector implements the GitHub connector
type Connector struct {
	API   API
	Owner string
	Repo  string
}

// Init takes the initialization of connector, e.g. reading environment vars etc
func (c *Connector) Init() {
	c.API = NewAPIClient()
	c.Owner = "artem-sidorenko"
	c.Repo = "chef-cups"
}

// GetTags returns the git tags
func (c *Connector) GetTags() (connectors.Tags, error) {
	tags, err := c.API.ListTags(c.Owner, c.Repo)
	if err != nil {
		return nil, err
	}

	var ret connectors.Tags
	for _, tag := range tags {
		commit, err := c.API.GetCommit(c.Owner, c.Repo, tag.Commit.GetSHA())

		if err != nil {
			return nil, err
		}

		ret = append(ret, connectors.Tag{
			Name:   tag.GetName(),
			Commit: commit.Commit.GetSHA(),
			Date:   commit.Commit.Committer.GetDate(),
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
