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

// Package generate implments the generate command
package generate

import (
	"io"
	"os"
	"time"

	"github.com/artem-sidorenko/chagen/commands"
	"github.com/artem-sidorenko/chagen/connectors"
	_ "github.com/artem-sidorenko/chagen/connectors/github" //enable github
	"github.com/artem-sidorenko/chagen/data"
	"github.com/artem-sidorenko/chagen/generator"

	"github.com/urfave/cli"
)

// Stdout references the Stdout writer for generate command
var Stdout io.Writer = os.Stdout // nolint: gochecknoglobals
// Connector references the connector ID used for generation
var Connector = "github" // nolint: gochecknoglobals

// Generate implements the CLI subcommand generate
func Generate(c *cli.Context) (err error) {
	tags, issues, mrs, err := getConnectorData(c.String("new-release"), c)
	if err != nil {
		return nil
	}

	gen := generator.New(data.NewReleases(tags, issues, mrs))

	// use stdout if - is given, otherwise create a new file
	filename := c.String("file")
	var wr io.Writer
	if filename != "-" {
		var file *os.File
		if file, err = os.Create(filename); err != nil {
			return err
		}

		defer func() {
			if cerr := file.Close(); err == nil && cerr != nil {
				err = cerr
			}
		}()

		wr = file
	} else {
		wr = Stdout
	}

	err = gen.Render(wr)

	return err
}

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

	connector, err = connectors.GetConnector(Connector)
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

	issues, err = connector.GetIssues()
	if err != nil {
		return nil, nil, nil, err
	}

	mrs, err = connector.GetMRs()
	if err != nil {
		return nil, nil, nil, err
	}

	return tags, issues, mrs, nil
}

// CLIFlags returns the possible CLI flags for this command
func CLIFlags() []cli.Flag {
	return []cli.Flag{
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
}

func init() { // nolint: gochecknoinits
	flags := CLIFlags()

	connectorFlags, _ := connectors.GetCLIFlags(Connector) // nolint: gosec

	flags = append(flags, connectorFlags...)

	commands.RegisterCommand(cli.Command{
		Name:      "generate",
		Usage:     "Generate a changelog",
		ArgsUsage: " ", // we do not have any args (only flags), so avoid this help message
		Flags:     flags,
		Action: func(c *cli.Context) error {
			if err := Generate(c); err != nil { // exit 1 and error message if we get any error reported
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	})
}
