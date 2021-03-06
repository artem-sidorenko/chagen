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

// Package connectors contains different connectors for fetching the data
package connectors

import (
	"fmt"

	"github.com/urfave/cli"
)

// NewConnectorFunc describes the constructor of Connector
type NewConnectorFunc func(*cli.Context) (Connector, error)

// ConnectorCLIFlags describes the function, which returns the configured
// CLI flags for particular connector
type ConnectorCLIFlags func() []cli.Flag

type connector struct {
	name         string
	newConnector NewConnectorFunc
	CLIFlags     ConnectorCLIFlags
}

var connectors = make(map[string]connector) // nolint: gochecknoglobals

// RegisterConnector registers the new connector for fetching the data.
// id is used as internal id or as value for CLI flag
// name is a text name of connector for humans (e.g. help pages)
// newConnector is the connector constructor function
// CLIFlag is a function, which returns the configured CLI flags
func RegisterConnector(id, name string, newConnector NewConnectorFunc, CLIFlags ConnectorCLIFlags) {
	connectors[id] = connector{
		name:         name,
		newConnector: newConnector,
		CLIFlags:     CLIFlags,
	}
}

// NewConnector returns the Connector of given id
// if this connector is missing, error is returned
func NewConnector(id string, ctx *cli.Context) (Connector, error) {
	if err := checkConnector(id); err != nil {
		return nil, err
	}

	conn, err := connectors[id].newConnector(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// CLIFlags returns the registered CLI flags of given connector
// if this connector is missing, error is returned
func CLIFlags(id string) ([]cli.Flag, error) {
	if err := checkConnector(id); err != nil {
		return nil, err
	}
	return connectors[id].CLIFlags(), nil
}

// checkConnector checks if given connector is registered.
// Returns nil if everything ok, error otherwise
func checkConnector(id string) error {
	if ConnectorRegistered(id) {
		return nil
	}
	return fmt.Errorf("unknown connector: %s", id)
}

// ConnectorRegistered verifies if given connector is registered
// returns true if yes, false otherwise
func ConnectorRegistered(id string) bool {
	_, ok := connectors[id]
	return ok
}

// ResetConnectors resets all registered connectors
// this is mainly used in the tests
func ResetConnectors() {
	connectors = make(map[string]connector)
}

// RegisteredConnectors returns a slice of registered connector ids
func RegisteredConnectors() []string {
	var ret []string
	for k := range connectors {
		ret = append(ret, k)
	}
	return ret
}
