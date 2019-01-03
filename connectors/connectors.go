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

	"github.com/artem-sidorenko/chagen/data"
	"github.com/urfave/cli"
)

// Connector describes the interface of connectors to the data sources
type Connector interface {
	Init(*cli.Context) error
	GetTags() (data.Tags, error)
	GetIssues() (data.Issues, error)
	GetMRs() (data.MRs, error)
	GetNewTagURL(string) (string, error)
}

type connector struct {
	name      string
	connector Connector
	flags     []cli.Flag
}

var connectors = make(map[string]connector) // nolint: gochecknoglobals

// RegisterConnector registers the new connector for fetching the data.
// id is used as internal id or as value for CLI flag
// name is a text name of connector for humans (e.g. help pages)
// c is the connector interface
// f is a slice with CLI flags of this connector
func RegisterConnector(id, name string, c Connector, f []cli.Flag) {
	connectors[id] = connector{
		name:      name,
		connector: c,
		flags:     f,
	}
}

// GetConnector returns the Connector of given id
// if this connector is missing, error is returned
func GetConnector(id string) (Connector, error) {
	if err := checkConnector(id); err != nil {
		return nil, err
	}
	return connectors[id].connector, nil
}

// GetCLIFlags returns the registered CLI flags of given connector
// if this connector is missing, error is returned
func GetCLIFlags(id string) ([]cli.Flag, error) {
	if err := checkConnector(id); err != nil {
		return nil, err
	}
	return connectors[id].flags, nil
}

// checkConnector checks if given connector is registered.
// Returns nil if everything ok, error otherwise
func checkConnector(id string) error {
	if _, ok := connectors[id]; !ok {
		return fmt.Errorf("Unknown connector: %s", id)
	}
	return nil
}
