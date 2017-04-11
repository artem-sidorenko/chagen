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
	"os"

	"github.com/artem-sidorenko/chagen/action"
	"github.com/urfave/cli"
)

var version = "0.0.1-dev"

const usage = "Changelog generator for your projects"

func main() {
	app := cli.NewApp()
	app.Name = "chagen"
	app.Version = version
	app.Usage = usage
	app.ArgsUsage = " " // we do not have any args (only flags), so avoid this help message
	app.Commands = commands()
	app.Authors = []cli.Author{
		{
			Name:  "Artem Sidorenko",
			Email: "artem@posteo.de",
		},
	}
	app.Run(os.Args)
}

func commands() []cli.Command {
	return []cli.Command{
		{
			Name:      "generate",
			Usage:     "Generate a changelog",
			ArgsUsage: " ", // we do not have any args (only flags), so avoid this help message
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Usage: "File name of changelog, - is accepted for stdout",
					Value: "CHANGELOG.md",
				},
			},
			Action: func(c *cli.Context) error {
				err := action.Generate(c.String("file"))
				if err != nil { // exit 1 and error message if we get any error reported
					return cli.NewExitError(err, 1)
				}
				return nil
			},
		},
	}
}
