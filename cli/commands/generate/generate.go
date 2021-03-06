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
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/artem-sidorenko/chagen/cli/commands"
	"github.com/artem-sidorenko/chagen/data"
	"github.com/artem-sidorenko/chagen/generator"
	"github.com/artem-sidorenko/chagen/internal/output"

	// some of connectors control functionality will be moved to source at some point
	_ "github.com/artem-sidorenko/chagen/datasource"
	"github.com/artem-sidorenko/chagen/datasource/connectors"

	"github.com/urfave/cli"
)

// Stdout references the Stdout writer for generate command
var Stdout io.Writer = output.Stdout // nolint: gochecknoglobals
// ProgressWriter references the writer for progress information
var ProgressWriter io.Writer = output.Stderr // nolint: gochecknoglobals

// Generate implements the CLI subcommand generate
func Generate(ctx *cli.Context) error { // nolint: gocyclo
	var filterRe *regexp.Regexp
	var excludeLabels []string
	var err error

	// verify the given endpoint, it should be one of supported connectors
	if ctx.String("endpoint") == "" {
		return fmt.Errorf("endpoint type is missing")
	}
	connector := ctx.String("endpoint")
	if !connectors.ConnectorRegistered(connector) {
		return fmt.Errorf("given endpoint isn't supported: %v", connector)
	}

	if !ctx.Bool("no-filter-tags") { // if the flag is not there, lets apply the filter
		filterReStr := ctx.String("filter-tags")
		if filterReStr == "" {
			return fmt.Errorf("regular expression for tag filtering should be defined")
		}
		if filterRe, err = regexp.Compile(filterReStr); err != nil {
			return fmt.Errorf("can't compile the regular expression: %v", err)
		}
	}

	ls := ctx.String("exclude-labels")
	if ls != "" {
		excludeLabels = strings.Split(ls, ",")
		//trim spaces from label names
		for i := range excludeLabels {
			excludeLabels[i] = strings.Trim(excludeLabels[i], " ")
		}
	}

	conn, err := connectors.NewConnector(connector, ctx)
	if err != nil {
		return err
	}

	exists, err := conn.RepositoryExists()
	if err != nil {
		return err
	}

	if !exists {
		// TODO: this should provide detailed information about repository: owner, repo name
		return fmt.Errorf("project not found")
	}

	tags, issues, mrs, err := getConnectorData(
		conn,
		filterRe,
		excludeLabels,
		ctx.String("new-release"),
	)
	if err != nil {
		return err
	}

	gen := generator.New(data.NewReleases(tags, issues, mrs))

	// use stdout if - is given, otherwise create a new file
	filename := ctx.String("file")
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

// collectData fans-in data from different channels to the data structures
func collectData( // nolint: gocyclo
	ctx context.Context,
	ctags <-chan data.Tag,
	cissues <-chan data.Issue,
	cmrs <-chan data.MR,
	cerr <-chan error,
) (
	data.Tags,
	data.Issues,
	data.MRs,
	error,
) {
	var (
		tags   data.Tags
		issues data.Issues
		mrs    data.MRs
	)

	for {
		select {
		case <-ctx.Done():
			return tags, issues, mrs, ctx.Err()
		case err, ok := <-cerr:
			if ok {
				return nil, nil, nil, err
			}
		case t, ok := <-ctags:
			if ok {
				tags = append(tags, t)
			} else { // tags are finished, nil the channel
				ctags = nil
			}
		case i, ok := <-cissues:
			if ok {
				issues = append(issues, i)
			} else { // issues are finished, nil the channel
				cissues = nil
			}
		case m, ok := <-cmrs:
			if ok {
				mrs = append(mrs, m)
			} else { // MRs are finished, nil the channel
				cmrs = nil
			}
		}
		// all channels finished, return data
		if ctags == nil && cissues == nil && cmrs == nil {
			return tags, issues, mrs, nil
		}
	}
}

// getConnectorData returns all needed data from connector
// if newRelease is specified, a new releases for
// untagged activities is created
func getConnectorData(
	conn connectors.Connector,
	tagsFilter *regexp.Regexp,
	excludeLabels []string,
	newRelease string,
) (data.Tags, data.Issues, data.MRs, error) {

	var (
		tags   data.Tags
		issues data.Issues
		mrs    data.MRs
	)

	// one minute for data collection should be enougth for now
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// we use cerr to track the possible errors in all goroutines invoked here
	cerr := make(chan error)

	ctags, ctagscounter, cmaxtags := conn.Tags(ctx, cerr)
	cissues, cissuescounter, cmaxissues := conn.Issues(ctx, cerr)
	cmrs, cmrscounter, cmaxmrs := conn.MRs(ctx, cerr)

	// invoke the progress printer
	printProgress(
		ctx, ProgressWriter,
		ctagscounter, cmaxtags,
		cissuescounter, cmaxissues,
		cmrscounter, cmaxmrs)

	//fan-in everything
	tags, issues, mrs, err := collectData(
		ctx,
		ctags,
		cissues,
		cmrs,
		cerr,
	)
	// release the err channel and only then process the possible errors
	close(cerr)
	if err != nil {
		return nil, nil, nil, err
	}

	// we should apply the filter to the tags
	if tagsFilter != nil {
		tags = data.FilterTags(tags, tagsFilter)
	}

	if newRelease != "" {
		var relURL string
		relURL, err = conn.GetNewTagURL(newRelease)
		if err != nil {
			return nil, nil, nil, err
		}

		tags = append(tags, data.Tag{
			Name: newRelease,
			Date: time.Now(),
			URL:  relURL,
		})
	}

	// we should filter the labels
	if len(excludeLabels) > 0 {
		issues = data.FilterIssuesByLabel(issues, excludeLabels)
		mrs = data.FilterMRsByLabel(mrs, excludeLabels)
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
			Usage: "Use the given release name and create a new release for all changes after the last tagged release", // nolint: lll
		},
		cli.StringFlag{
			Name:  "filter-tags, t",
			Usage: "Only use tags, which match to the given regular expression",
			Value: `^v\d+\.\d+\.\d+$`,
		},
		cli.BoolFlag{
			Name:  "no-filter-tags",
			Usage: "Disable filtering of tags",
		},
		cli.StringFlag{
			Name:  "exclude-labels",
			Usage: "Exclude issues and MRs/PRs with specified labels `x,y,z`",
			Value: "duplicate, question, invalid, wontfix, no changelog",
		},
		cli.StringFlag{
			Name:  "endpoint",
			Usage: "API endpoint type: " + strings.Join(connectors.RegisteredConnectors(), ", "),
			Value: "github",
		},
	}
}

func init() { // nolint: gochecknoinits
	flags := CLIFlags()

	for _, conn := range connectors.RegisteredConnectors() {
		connectorFlags, err := connectors.CLIFlags(conn)
		if err != nil {
			panic(err)
		}

		flags = append(flags, connectorFlags...)
	}

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
