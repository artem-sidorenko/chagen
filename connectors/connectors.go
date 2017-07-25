// Package connectors contains different connectors for fetching the data
package connectors

import (
	"fmt"

	"github.com/artem-sidorenko/chagen/data"
)

// Connector describes the interface of connectors to the data sources
type Connector interface {
	Init()
	GetTags() (data.Tags, error)
	GetIssues() (data.Issues, error)
	GetMRs() (data.MRs, error)
}

type connector struct {
	name      string
	connector Connector
}

var connectors = make(map[string]connector)

// RegisterConnector registers the new connector for fetching the data.
// id is used as internal id or as value for CLI flag
// name is a text name of connector for humans (e.g. help pages)
// other parameter is the connector interface
func RegisterConnector(id, name string, c Connector) {
	connectors[id] = connector{
		name:      name,
		connector: c,
	}
}

// GetConnector returns the Connector of given id
// if this connector is missing, error is returned
func GetConnector(id string) (Connector, error) {
	if _, ok := connectors[id]; !ok {
		return nil, fmt.Errorf("Unknown connector: %s", id)
	}
	return connectors[id].connector, nil
}
