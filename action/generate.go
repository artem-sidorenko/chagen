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

// Package action contains the specific implementations for subcommands
package action

import (
	"os"

	"github.com/artem-sidorenko/chagen/generator"
	"github.com/artem-sidorenko/chagen/internal/testdata"
)

// Generate implements the CLI subcommand generate
func Generate(filename string) error {
	gen := generator.Generator{
		Releases: testdata.GeneratorDataStructure, //for now we use testdata as source
	}

	// use stdout if - is given, otherwise create a new file
	wr := os.Stdout
	if filename != "-" {
		var err error
		wr, err = os.Create(filename)
		if err != nil {
			return err
		}
		defer wr.Close()
	}

	err := gen.Render(wr)
	if err != nil {
		return err
	}

	return nil
}
