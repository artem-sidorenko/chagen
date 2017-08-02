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

package commands

import (
	"os"
	"time"

	"github.com/artem-sidorenko/chagen/connectors"
	_ "github.com/artem-sidorenko/chagen/connectors/github" //enable github
	"github.com/artem-sidorenko/chagen/data"
	"github.com/artem-sidorenko/chagen/generator"
	"github.com/urfave/cli"
)

// Generate implements the CLI subcommand generate
func Generate(c *cli.Context) (err error) { // nolint: gocyclo
	connector, err := connectors.GetConnector("github")
	if err != nil {
		return
	}

	err = connector.Init(c)
	if err != nil {
		return
	}

	tags, err := connector.GetTags()
	if err != nil {
		return
	}
	if rel := c.String("new-release"); rel != "" {
		var relURL string
		relURL, err = connector.GetNewTagURL(rel)
		if err != nil {
			return
		}

		tags = append(tags, data.Tag{
			Name: rel,
			Date: time.Now(),
			URL:  relURL,
		})
	}
	tags.Sort()

	issues, err := connector.GetIssues()
	if err != nil {
		return
	}
	issues.Sort()

	mrs, err := connector.GetMRs()
	if err != nil {
		return
	}
	mrs.Sort()

	releases := data.NewReleases(tags, issues, mrs)

	gen := generator.Generator{
		Releases: releases,
	}

	// use stdout if - is given, otherwise create a new file
	filename := c.String("file")
	wr := os.Stdout
	if filename != "-" {
		wr, err = os.Create(filename)
		if err != nil {
			return
		}
		defer func() {
			if cerr := wr.Close(); err == nil && cerr != nil {
				err = cerr
			}
		}()
	}

	err = gen.Render(wr)

	return
}

func init() {
	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "File name of changelog, - is accepted for stdout",
			Value: "CHANGELOG.md",
		},
		cli.StringFlag{
			Name:  "new-release, r",
			Usage: "Create a new release for all issues and changes after the last release",
		},
	}

	connectorFlags, _ := connectors.GetCLIFlags("github")

	flags = append(flags, connectorFlags...)

	RegisterCommand(cli.Command{
		Name:      "generate",
		Usage:     "Generate a changelog",
		ArgsUsage: " ", // we do not have any args (only flags), so avoid this help message
		Flags:     flags,
		Action: func(c *cli.Context) error {
			err := Generate(c)
			if err != nil { // exit 1 and error message if we get any error reported
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	})
}
