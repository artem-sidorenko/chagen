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

package main

import (
	"fmt"
	"os"

	"github.com/artem-sidorenko/chagen/commands"
	hcli "github.com/artem-sidorenko/chagen/helpers/cli"

	"github.com/urfave/cli"
)

var version = "unknown" // nolint: gochecknoglobals

const usage = "Changelog generator for your projects"

func main() {
	app := cli.NewApp()
	app.Version = version
	app.OnUsageError = hcli.OnUsageError
	app.Usage = usage
	// we do not have any args (only flags), so avoid this help message
	app.ArgsUsage = " "
	app.Commands = commands.GetCommands()
	app.Authors = []cli.Author{
		{
			Name:  "Artem Sidorenko",
			Email: "artem@posteo.de",
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}
