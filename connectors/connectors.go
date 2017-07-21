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
	"time"
)

// Tag describes a git tag
type Tag struct {
	Name   string
	Commit string
	Date   time.Time
	URL    string
}

// Tags is a slice with Tag elements
type Tags []Tag

// Issue describes an issue in the bug tracker
type Issue struct {
	ID         int
	Name       string
	ClosedDate time.Time
	URL        string
}

// Issues is a slice with Issue elements
type Issues []Issue

// MR describes a Pull or Merge Request
type MR struct {
	ID         int
	Name       string
	MergedDate time.Time
}

// MRs is a slice with MR elements
type MRs []MR

// Connector describes the interface of connectors to the data sources
type Connector interface {
	Init()
	GetTags() (Tags, error)
	GetIssues() (Issues, error)
	GetMRs() (MRs, error)
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
