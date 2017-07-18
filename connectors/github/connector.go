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
func (c *Connector) GetTags() (ret connectors.Tags) {
	tags := c.API.ListTags(c.Owner, c.Repo)
	for _, tag := range tags {
		commit := c.API.GetCommit(c.Owner, c.Repo, tag.Commit.GetSHA())

		ret = append(ret, connectors.Tag{
			Name:   tag.GetName(),
			Commit: commit.Commit.GetSHA(),
			Date:   commit.Commit.Committer.GetDate(),
		})
	}
	return
}

func init() {
	connectors.RegisterConnector("github", "GitHub", &Connector{})
}
