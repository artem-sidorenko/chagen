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

	"github.com/artem-sidorenko/chagen/generator"
	"github.com/artem-sidorenko/chagen/internal/testdata"
	"github.com/urfave/cli"
)

// Generate implements the CLI subcommand generate
func Generate(filename string) (err error) {
	gen := generator.Generator{
		Releases: testdata.GeneratorDataStructure, //for now we use testdata as source
	}

	// use stdout if - is given, otherwise create a new file
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
	RegisterCommand(cli.Command{
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
			err := Generate(c.String("file"))
			if err != nil { // exit 1 and error message if we get any error reported
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	})
}
