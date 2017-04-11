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

package generator

import (
	"bytes"
	"testing"

	"github.com/artem-sidorenko/chagen/data"
	"github.com/artem-sidorenko/chagen/internal/testdata"
)

func TestGenerator_Render(t *testing.T) {
	type fields struct {
		Releases []data.Release
	}
	tests := []struct {
		name    string
		fields  fields
		wantWr  string
		wantErr bool
	}{
		{
			name: "verify testdata structure",
			fields: fields{
				Releases: testdata.GeneratorDataStructure,
			},
			wantWr:  testdata.GeneratedChangelog,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				Releases: tt.fields.Releases,
			}
			wr := &bytes.Buffer{}
			if err := g.Render(wr); (err != nil) != tt.wantErr {
				t.Errorf("Generator.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWr := wr.String(); gotWr != tt.wantWr {
				t.Errorf("Generator.Render() = %v, want %v", gotWr, tt.wantWr)
			}
		})
	}
}
