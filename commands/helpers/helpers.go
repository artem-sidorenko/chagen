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

// Package helpers contains some helper functions,
// which are used by several commands
package helpers

import (
	"time"

	"github.com/artem-sidorenko/chagen/connectors"
	"github.com/artem-sidorenko/chagen/data"
	"github.com/artem-sidorenko/chagen/generator"

	"github.com/urfave/cli"
)

// getConnectorData returns all needed data from connector
// if newRelease is specified, a new releases for
// untagged activities is created
func getConnectorData(newRelease string, c *cli.Context) (data.Tags, data.Issues, data.MRs, error) {
	var (
		connector connectors.Connector
		tags      data.Tags
		issues    data.Issues
		mrs       data.MRs
		err       error
	)

	connector, err = connectors.GetConnector("github")
	if err != nil {
		return nil, nil, nil, err
	}

	err = connector.Init(c)
	if err != nil {
		return nil, nil, nil, err
	}

	tags, err = connector.GetTags()
	if err != nil {
		return nil, nil, nil, err
	}

	if newRelease != "" {
		var relURL string
		relURL, err = connector.GetNewTagURL(newRelease)
		if err != nil {
			return nil, nil, nil, err
		}

		tags = append(tags, data.Tag{
			Name: newRelease,
			Date: time.Now(),
			URL:  relURL,
		})
	}
	tags.Sort()

	issues, err = connector.GetIssues()
	if err != nil {
		return nil, nil, nil, err
	}
	issues.Sort()

	mrs, err = connector.GetMRs()
	if err != nil {
		return nil, nil, nil, err
	}
	mrs.Sort()

	return tags, issues, mrs, nil
}

// NewGenerator returns a new generator,
// which is filled and initialized with data
func NewGenerator(c *cli.Context) (*generator.Generator, error) {
	tags, issues, mrs, err := getConnectorData(c.String("new-release"), c)
	if err != nil {
		return nil, err
	}

	gen := generator.Generator{
		Releases: data.NewReleases(tags, issues, mrs),
	}

	return &gen, nil
}
