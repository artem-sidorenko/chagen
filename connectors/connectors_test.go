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
	"errors"
	"reflect"
	"testing"

	"github.com/artem-sidorenko/chagen/connectors"
)

type testConnector struct{}

func (t *testConnector) Init()                                 {}
func (t *testConnector) GetTags() (connectors.Tags, error)     { return nil, nil }
func (t *testConnector) GetIssues() (connectors.Issues, error) { return nil, nil }
func (t *testConnector) GetMRs() (connectors.MRs, error)       { return nil, nil }

func TestGetConnector(t *testing.T) {
	connectors.RegisterConnector("testexisting", "TestExisting", &testConnector{})

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
			wantErr: errors.New("Unknown connector: testmissing"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := connectors.GetConnector(tt.args.id)
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("GetConnector() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConnector() = %v, want %v", got, tt.want)
			}
		})
	}
}