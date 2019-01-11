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

// Package cli provides some helper functions for github.com/urfave/cli
package cli

import (
	"fmt"

	"github.com/urfave/cli"
)

// OnUsageError represents a workaround for https://github.com/urfave/cli/issues/610
// once there is a fix upstream, this can be removed
func OnUsageError(context *cli.Context, err error, isSubcommand bool) error {
	fmt.Fprintf( // nolint: errcheck
		cli.ErrWriter, "%s - %s %s\n\n",
		context.App.Name, "incorrect usage:", err.Error(),
	)
	cli.ShowAppHelp(context) // nolint: gosec, errcheck
	return cli.NewExitError("", 1)
}
}
