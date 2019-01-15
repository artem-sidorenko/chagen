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

// Package cli provides the CLI functionality and controls github.com/urfave/cli
package cli

import (
	"fmt"
	"os"

	"github.com/artem-sidorenko/chagen/cli/commands"
	_ "github.com/artem-sidorenko/chagen/cli/commands/generate" // enable generate subcommand
	"github.com/artem-sidorenko/chagen/internal/info"

	"github.com/urfave/cli"
)

const usage = "Changelog generator for your projects"

// Run provides the main encapsulation of control logic for CLI
// basically its like main()
func Run() {
	app := cli.NewApp()
	app.Version = info.Version()
	app.OnUsageError = onUsageError
	app.ExitErrHandler = exitErrHandler
	app.Usage = usage
	// we do not have any args (only flags), so avoid this help message
	app.ArgsUsage = " "
	app.Commands = commands.GetCommands()
	app.Authors = []cli.Author{
		{
			Name:  info.Author,
			Email: info.Email,
		},
	}
	err := app.Run(os.Args)
	// Usually this should not happen and err should be catched within app.Run
	// because of our own ExitErrHandler. Lets have it here just for the case
	// with a different exit code
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err) // nolint: errcheck
		os.Exit(10)
	}
}

// onUsageError represents a workaround for https://github.com/urfave/cli/issues/610
// once there is a fix upstream, this can be removed
func onUsageError(context *cli.Context, err error, isSubcommand bool) error {
	fmt.Fprintf( // nolint: errcheck
		cli.ErrWriter, "%s - %s %s\n\n",
		context.App.Name, "incorrect usage:", err.Error(),
	)
	cli.ShowAppHelp(context) // nolint: gosec, errcheck
	return cli.NewExitError("", 1)
}

// exitErrHandler implements cli.ExitErrHandlerFunc
// we make it simple, we always return exit code 1
func exitErrHandler(_ *cli.Context, err error) {
	if err.Error() != "" {
		fmt.Fprintf(cli.ErrWriter, "Error: %+v\n", err) // nolint: errcheck
	}
	cli.OsExiter(1)
}
