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

// Package client implements the access to the github API package via
// interface and own wrapper client
package client

import (
	"context"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// NewGitHubClient intialized and returns a new gitHubClient
// Uses AccessToken for oauth2 authentication if not empty
func NewGitHubClient(ctx context.Context, AccessToken string) *GitHubClient {
	var tc *http.Client

	if AccessToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: AccessToken},
		)
		tc = oauth2.NewClient(ctx, ts)
	}

	client := github.NewClient(tc)

	return &GitHubClient{
		Repositories: client.Repositories,
		Issues:       client.Issues,
		PullRequests: client.PullRequests,
	}
}
