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
	"os"

	"github.com/artem-sidorenko/chagen/source/connectors"
	"github.com/artem-sidorenko/chagen/source/connectors/github/internal/client"

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
