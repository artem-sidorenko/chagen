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

package connectors_test

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"testing"

	"github.com/artem-sidorenko/chagen/data"
	"github.com/artem-sidorenko/chagen/datasource/connectors"

	"github.com/urfave/cli"
)

type testConnector struct{}

func (t *testConnector) Tags(_ context.Context, _ chan<- error) (
	<-chan data.Tag, <-chan bool, <-chan int,
) {
	return nil, nil, nil
}
func (t *testConnector) Issues(_ context.Context, _ chan<- error) (
	<-chan data.Issue, <-chan bool, <-chan int,
) {
	return nil, nil, nil
}
func (t *testConnector) MRs(_ context.Context, _ chan<- error) (
	<-chan data.MR, <-chan bool, <-chan int,
) {
	return nil, nil, nil
}
func (t *testConnector) GetNewTagURL(string) (string, error) { return "", nil }
func (t *testConnector) RepositoryExists() (bool, error)     { return true, nil }

func newTestConnector(_ *cli.Context) (connectors.Connector, error) {
	return &testConnector{}, nil
}

func registerConnectors(ids []string) {
	connectors.ResetConnectors()
	for _, conn := range ids {
		connectors.RegisterConnector(conn, conn, newTestConnector, nil)
	}
}

func CLIFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "test, t",
			Usage: "Testing",
		},
	}
}

func TestCLIFlags(t *testing.T) {
	connectors.RegisterConnector("testexisting", "TestExisting", newTestConnector, CLIFlags)

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    []cli.Flag
		wantErr error
	}{
		{
			name: "Connector exists",
			args: args{
				id: "testexisting",
			},
			want: []cli.Flag{
				cli.BoolFlag{
					Name:  "test, t",
					Usage: "Testing",
				},
			},
		},
		{
			name: "Connector does not exist",
			args: args{
				id: "testmissing",
			},
			wantErr: errors.New("unknown connector: testmissing"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := connectors.CLIFlags(tt.args.id)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("CLIFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CLIFlags() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestNewConnector(t *testing.T) {
	connectors.ResetConnectors()
	connectors.RegisterConnector("testexisting", "TestExisting", newTestConnector, nil)

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    connectors.Connector
		wantErr error
	}{
		{
			name: "Connector exists",
			args: args{
				id: "testexisting",
			},
			want: &testConnector{},
		},
		{
			name: "Connector does not exist",
			args: args{
				id: "testmissing",
			},
			wantErr: errors.New("unknown connector: testmissing"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := connectors.NewConnector(tt.args.id, nil)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("NewConnector() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConnector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConnectorRegistered(t *testing.T) {
	tests := []struct {
		name          string
		regConnectors []string
		id            string
		want          bool
	}{
		{
			name:          "Registered connector is requested",
			regConnectors: []string{"testconn1", "testconn2"},
			id:            "testconn2",
			want:          true,
		},
		{
			name:          "Not registered connector is requested",
			regConnectors: []string{"testconn1"},
			id:            "testconn2",
			want:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registerConnectors(tt.regConnectors)

			if got := connectors.ConnectorRegistered(tt.id); got != tt.want {
				t.Errorf("ConnectorRegistered() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisteredConnectors(t *testing.T) {
	tests := []struct {
		name          string
		regConnectors []string
		want          []string
	}{
		{
			name:          "Registered connectors",
			regConnectors: []string{"testconn1", "testconn2"},
			want:          []string{"testconn1", "testconn2"},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registerConnectors(tt.regConnectors)

			got := connectors.RegisteredConnectors()
			sort.Strings(got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisteredConnectors() = %v, want %v", got, tt.want)
			}
		})
	}
}
