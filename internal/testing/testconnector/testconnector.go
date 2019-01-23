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
	"context"
	"time"

	"github.com/artem-sidorenko/chagen/connectors"
	"github.com/artem-sidorenko/chagen/data"
	"github.com/artem-sidorenko/chagen/internal/testing/testconnector/testdata"

	"github.com/urfave/cli"
)

// nolint: gochecknoglobals
var (
	// RetTestingTag controls whenether the testconnector should return the tag testingtag
	RetTestingTag = false
	// RepositoryExistsFail controls whenether the testconnector should
	// fail in the RepositoryExists() call
	RepositoryExistsFail = false
)

// Connector implements the test connector
type Connector struct {
}

// RepositoryExists checks if referenced repository is present
func (c *Connector) RepositoryExists() (bool, error) {
	return !RepositoryExistsFail, nil
}

// Tags implements the connectors.Connector interface
func (c *Connector) Tags(
	_ context.Context,
	cerr chan<- error,
) (
	<-chan data.Tag,
	<-chan int,
) {
	tags := testdata.Tags()

	if RetTestingTag {
		tags = append(tags, data.Tag{
			Name:   "testingtag",
			Date:   time.Unix(1147783647, 0),
			Commit: "ad59c6b54ba53f54383d7c4661bdd4e29fe87c15",
			URL:    "https://test.example.com/tags/testingtag",
		})
	}

	ctags := make(chan data.Tag)

	go func() {
		defer close(ctags)

		for _, t := range tags {
			ctags <- t
		}
	}()

	return ctags, nil
}

// Issues implements the connectors.Connector interface
func (c *Connector) Issues(
	_ context.Context,
	cerr chan<- error,
) (
	<-chan data.Issue,
	<-chan int,
) {
	cissues := make(chan data.Issue)

	go func() {
		defer close(cissues)

		for _, t := range testdata.Issues() {
			cissues <- t
		}
	}()

	return cissues, nil
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
		{
			ID:         4,
			Name:       "MR 4",
			MergedDate: time.Unix(1299483647, 0),
			Author:     "testauthor",
			AuthorURL:  "https://test.example.com/authors/testauthor",
			URL:        "https://test.example.com/mrs/4",
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

// CLIFlags describes the flags of connector
func CLIFlags() []cli.Flag {
	return []cli.Flag{}
}

func init() { // nolint: gochecknoinits
	connectors.RegisterConnector("testconnector", "TestConnector", New, CLIFlags)
}
