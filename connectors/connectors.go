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

var connectors = make(map[string]connector)

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
