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
	"errors"
	"html/template"
	"testing"
	"time"

	"github.com/artem-sidorenko/chagen/cli/commands/generate"
	tcli "github.com/artem-sidorenko/chagen/internal/testing/cli"
	"github.com/artem-sidorenko/chagen/internal/testing/testconnector"

	_ "github.com/artem-sidorenko/chagen/internal/testing/testconnector"
)

func genOutput(newRelease, testingTag, secondTag, excludedIssue bool) string {
	// nolint: lll
	tpl := `Changelog
=========

{{- if .NewRelease }}

## [v10.10.0](http://test.example.com/releases/v10.10.0) ({{.NewReleaseDate}})

Closed issues
-------------
- Issue 3 [\#3](http://test.example.com/issues/3)
{{- if .ExcludedIssue }}
- Issue 6 [\#6](http://test.example.com/issues/6)
{{- else }}
- Issue 5 [\#5](http://test.example.com/issues/5)
{{- end }}
- Issue 4 [\#4](http://test.example.com/issues/4)

Merged pull requests
--------------------
- MR 4 [\#4](https://test.example.com/mrs/4) ([testauthor](https://test.example.com/authors/testauthor))

{{- end }}

## [v0.0.3](https://test.example.com/tags/v0.0.3) (13.07.2009)

Closed issues
-------------
- Issue 2 [\#2](http://test.example.com/issues/2)

Merged pull requests
--------------------
- MR 2 [\#2](https://test.example.com/mrs/2) ([testauthor](https://test.example.com/authors/testauthor))

{{- if not .SecondTag }}
- MR 3 [\#3](https://test.example.com/mrs/3) ([testauthor](https://test.example.com/authors/testauthor))
{{- end }}

{{- if .TestingTag }}

## [testingtag](https://test.example.com/tags/testingtag) (16.05.2006)

{{- end }}

{{- if .SecondTag }}

## [v0.0.2](https://test.example.com/tags/v0.0.2) (13.05.2006)

Merged pull requests
--------------------
- MR 3 [\#3](https://test.example.com/mrs/3) ([testauthor](https://test.example.com/authors/testauthor))

{{- end }}

## [v0.0.1](https://test.example.com/tags/v0.0.1) (12.03.2003)

Closed issues
-------------
- Issue 1 [\#1](http://test.example.com/issues/1)

Merged pull requests
--------------------
- MR 1 [\#1](https://test.example.com/mrs/1) ([testauthor](https://test.example.com/authors/testauthor))

*This Changelog was automatically generated with [chagen unknown](https://github.com/artem-sidorenko/chagen)*`

	input := struct {
		NewRelease     bool
		NewReleaseDate string
		TestingTag     bool
		SecondTag      bool
		ExcludedIssue  bool
	}{newRelease, time.Now().Format("02.01.2006"), testingTag, secondTag, excludedIssue}

	t := template.Must(template.New("Output template").Parse(tpl))

	buf := &bytes.Buffer{}

	t.Execute(buf, input)

	return buf.String()
}

func TestGenerate(t *testing.T) { // nolint: gocyclo
	type cliParams struct {
		newRelease    string
		noFilterTags  bool
		filterExpr    string
		excludeLabels string
	}

	tests := []struct {
		name       string
		cliParams  cliParams
		wantErr    error
		wantOutput string
	}{
		{
			name:       "Default flags",
			wantOutput: genOutput(false, false, true, false),
		},
		{
			name: "With new release flag",
			cliParams: cliParams{
				newRelease: "v10.10.0",
			},
			wantOutput: genOutput(true, false, true, false),
		},
		{
			name: "With --no-filter-tags",
			cliParams: cliParams{
				noFilterTags: true,
			},
			wantOutput: genOutput(false, true, true, false),
		},
		{
			name: "With customized filter",
			cliParams: cliParams{
				filterExpr: `^v\d+\.\d+\.(1|3)+$`,
			},
			wantOutput: genOutput(false, false, false, false),
		},
		{
			name: "With broken filter",
			cliParams: cliParams{
				filterExpr: "(abdc",
			},
			wantErr: errors.New("Can't compile the regular expression: error parsing regexp: missing closing ): `(abdc`"), // nolint: lll
		},
		{
			name: "With customized labels",
			cliParams: cliParams{
				excludeLabels: "issue5",
				newRelease:    "v10.10.0",
			},
			wantOutput: genOutput(true, false, true, true),
		},
	}
	for _, tt := range tests {
		cliFlags := map[string]string{
			"file": "-",
		}
		if tt.cliParams.newRelease != "" {
			cliFlags["new-release"] = tt.cliParams.newRelease
		}
		if tt.cliParams.noFilterTags {
			cliFlags["no-filter-tags"] = "true"
		}
		if tt.cliParams.filterExpr != "" {
			cliFlags["filter-tags"] = tt.cliParams.filterExpr
		}
		if tt.cliParams.excludeLabels != "" {
			cliFlags["exclude-labels"] = tt.cliParams.excludeLabels
		}
		ctx := tcli.TestContext(generate.CLIFlags(), cliFlags)

		output := &bytes.Buffer{}
		generate.Stdout = output
		generate.Connector = "testconnector"
		testconnector.RetTestingTag = true

		t.Run(tt.name, func(t *testing.T) {
			err := generate.Generate(ctx)
			out := string(output.String())

			if (err != nil && tt.wantErr == nil) ||
				(err == nil && tt.wantErr != nil) ||
				((err != nil && tt.wantErr != nil) && (err.Error() != tt.wantErr.Error())) {

				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if out != tt.wantOutput {
				t.Errorf("Generate() output = %v, wantOutput %v", out, tt.wantOutput)
			}
		})
	}
}
