/*
   Copyright 2019 Artem Sidorenko <artem@posteo.de>

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

package generate_test

import (
	"bytes"
	"testing"

	"github.com/artem-sidorenko/chagen/cli/commands/generate"
	tcli "github.com/artem-sidorenko/chagen/internal/testing/cli"

	_ "github.com/artem-sidorenko/chagen/internal/testing/testconnector"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name       string
		wantErr    bool
		wantOutput string
	}{
		{
			name:    "proper data",
			wantErr: false,
			// nolint: lll
			wantOutput: `Changelog
=========

## [v0.0.3](https://test.example.com/tags/v0.0.3) (13.07.2009)

Closed issues
-------------
- Issue 2 [\#2](http://test.example.com/issues/2)

Merged pull requests
--------------------
- MR 2 [\#2](https://test.example.com/mrs/2) ([testauthor](https://test.example.com/authors/testauthor))

## [v0.0.2](https://test.example.com/tags/v0.0.2) (13.05.2006)

Merged pull requests
--------------------
- MR 3 [\#3](https://test.example.com/mrs/3) ([testauthor](https://test.example.com/authors/testauthor))

## [v0.0.1](https://test.example.com/tags/v0.0.1) (12.03.2003)

Closed issues
-------------
- Issue 1 [\#1](http://test.example.com/issues/1)

Merged pull requests
--------------------
- MR 1 [\#1](https://test.example.com/mrs/1) ([testauthor](https://test.example.com/authors/testauthor))

*This Changelog was automatically generated with [chagen unknown](https://github.com/artem-sidorenko/chagen)*`,
		},
	}
	for _, tt := range tests {
		ctx := tcli.TestContext(generate.CLIFlags(), map[string]string{
			"file": "-",
		})

		output := &bytes.Buffer{}
		generate.Stdout = output
		generate.Connector = "testconnector"

		t.Run(tt.name, func(t *testing.T) {
			err := generate.Generate(ctx)
			out := string(output.String())
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if out != tt.wantOutput {
				t.Errorf("Generate() output = %v, wantOutput %v", out, tt.wantOutput)
			}
		})
	}
}
