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

// Package commands contains the specific implementations for subcommands
package commands

import (
	"github.com/urfave/cli"
)

var commands cli.Commands // nolint: gochecknoglobals

// RegisterCommand registers the new subcommand
func RegisterCommand(c cli.Command) {
	commands = append(commands, c)
}

// GetCommands returns the registered commands
func GetCommands() cli.Commands {
	return commands
}
