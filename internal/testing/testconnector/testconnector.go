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

// Package testconnector provides a static connector for testing
// purposes. This connector always delivers the same data
package testconnector

import (
	"time"

	"github.com/artem-sidorenko/chagen/connectors"
	"github.com/artem-sidorenko/chagen/data"

	"github.com/urfave/cli"
)

// Connector implements the test connector
type Connector struct{}

// GetTags implements the connectors.Connector interface
func (*Connector) GetTags() (data.Tags, error) {
	return data.Tags{
		{
			Name:   "v0.0.2",
			Date:   time.Unix(1147483647, 0),
			Commit: "b6a735dcb420a82865abe8c194900e59f6af9dea",
			URL:    "https://test.example.com/tags/v0.0.2",
		},
		{
			Name:   "v0.0.1",
			Date:   time.Unix(1047483647, 0),
			Commit: "d85645cbe6288cce5e5d066f8c7864040266cce3",
			URL:    "https://test.example.com/tags/v0.0.1",
		},
		{
			Name:   "v0.0.3",
			Date:   time.Unix(1247483647, 0),
			Commit: "25362c337d524025bf98e978059bf9bcd2b56221",
			URL:    "https://test.example.com/tags/v0.0.3",
		},
	}, nil
}

// GetIssues implements the connectors.Connector interface
func (*Connector) GetIssues() (data.Issues, error) {
	return data.Issues{
		{
			ID:         2,
			Name:       "Issue 2",
			ClosedDate: time.Unix(1247483647, 0),
			URL:        "http://test.example.com/issues/2",
		},
		{
			ID:         1,
			Name:       "Issue 1",
			ClosedDate: time.Unix(1047483647, 0),
			URL:        "http://test.example.com/issues/1",
		},
		{
			ID:         3,
			Name:       "Issue 3",
			ClosedDate: time.Unix(1347483647, 0),
			URL:        "http://test.example.com/issues/3",
		},
	}, nil
}

// GetMRs implements the connectors.Connector interface
func (*Connector) GetMRs() (data.MRs, error) {
	return data.MRs{
		{
			ID:         2,
			Name:       "MR 2",
			MergedDate: time.Unix(1247483647, 0),
			Author:     "testauthor",
			AuthorURL:  "https://test.example.com/authors/testauthor",
			URL:        "https://test.example.com/mrs/2",
		},
		{
			ID:         1,
			Name:       "MR 1",
			MergedDate: time.Unix(1047483647, 0),
			Author:     "testauthor",
			AuthorURL:  "https://test.example.com/authors/testauthor",
			URL:        "https://test.example.com/mrs/1",
		},
		{
			ID:         3,
			Name:       "MR 3",
			MergedDate: time.Unix(1057483647, 0),
			Author:     "testauthor",
			AuthorURL:  "https://test.example.com/authors/testauthor",
			URL:        "https://test.example.com/mrs/3",
		},
	}, nil
}

// GetNewTagURL implements the connectors.Connector interface
func (*Connector) GetNewTagURL(TagName string) (string, error) {
	return "http://test.example.com/releases/" + TagName, nil
}

// New creates a new Connector
func New(ctx *cli.Context) (connectors.Connector, error) {
	return &Connector{}, nil
}

func init() { // nolint: gochecknoinits
	connectors.RegisterConnector("testconnector", "TestConnector", New, []cli.Flag{})
}
