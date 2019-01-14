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

// Package cli provides helpers for CLI testing
package cli

import (
	"flag"

	"github.com/urfave/cli"
)

// TestContext creates the simulation of CLI flag setting
// we need this for testing, unfotunelly useful functions within
// cli package are private, so we have to do it by ourself
func TestContext(flags []cli.Flag, set map[string]string) *cli.Context {
	flagset := flag.NewFlagSet("", flag.ContinueOnError)
	for _, x := range flags {
		x.Apply(flagset)
	}
	for k, v := range set {
		if err := flagset.Set(k, v); err != nil {
			panic(err)
		}
	}

	return cli.NewContext(nil, flagset, nil)
}
