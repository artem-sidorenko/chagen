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
	"testing"

	"github.com/artem-sidorenko/chagen/connectors"
	"github.com/artem-sidorenko/chagen/data"

	"github.com/urfave/cli"
)

type testConnector struct{}

func (t *testConnector) Tags(_ context.Context, _ chan<- error) (<-chan data.Tag, <-chan int) {
	return nil, nil
}
func (t *testConnector) Issues(_ context.Context, _ chan<- error) (<-chan data.Issue, <-chan int) {
	return nil, nil
}
func (t *testConnector) MRs(_ context.Context, _ chan<- error) (<-chan data.MR, <-chan int) {
	return nil, nil
}
func (t *testConnector) GetNewTagURL(string) (string, error) { return "", nil }
func (t *testConnector) RepositoryExists() (bool, error)     { return true, nil }

func NewTestConnector(_ *cli.Context) (connectors.Connector, error) {
	return &testConnector{}, nil
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
	connectors.RegisterConnector("testexisting", "TestExisting", NewTestConnector, CLIFlags)

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
	connectors.RegisterConnector("testexisting", "TestExisting", NewTestConnector, nil)

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
