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
	"reflect"
	"testing"
	"time"

	"github.com/artem-sidorenko/chagen/cli/commands/generate"
	tcli "github.com/artem-sidorenko/chagen/internal/testing/cli"
	"github.com/artem-sidorenko/chagen/internal/testing/testconnector"

	_ "github.com/artem-sidorenko/chagen/internal/testing/testconnector"
)

func genOutput(newRelease, testingTag, minorTags, excludedIssue bool) string {
	// nolint: lll
	tpl := `Changelog
=========

{{- if .NewRelease }}

## [v10.10.0](http://test.example.com/releases/v10.10.0) ({{.NewReleaseDate}})

Closed issues
-------------
- Test issue title 13 [\#1234](http://test.example.com/issues/1234)
{{- if not .ExcludedIssue }}
- Test issue title 12 [\#1224](http://test.example.com/issues/1224)
{{- end }}

Merged pull requests
--------------------
- Test PR title 14 [\#2344](https://test.example.com/mrs/2344) ([te77st-user](https://test.example.com/authors/te77st-user))
- Test PR title 13 [\#2334](https://test.example.com/mrs/2334) ([test-user](https://test.example.com/authors/test-user))

{{- end }}

{{- if .SecondTag }}

## [v0.1.2](https://test.example.com/tags/v0.1.2) (20.03.2003)

## [v0.1.1](https://test.example.com/tags/v0.1.1) (19.03.2003)

Merged pull requests
--------------------
- Test PR title 10 [\#2304](https://test.example.com/mrs/2304) ([test-user](https://test.example.com/authors/test-user))

## [v0.1.0](https://test.example.com/tags/v0.1.0) (18.03.2003)

Closed issues
-------------
- Test issue title 9 [\#1294](http://test.example.com/issues/1294)

Merged pull requests
--------------------
- Test PR title 9 [\#2294](https://test.example.com/mrs/2294) ([test-user](https://test.example.com/authors/test-user))

{{- end }}

## [v0.0.9](https://test.example.com/tags/v0.0.9) (17.03.2003)

## [v0.0.8](https://test.example.com/tags/v0.0.8) (16.03.2003)

Merged pull requests
--------------------
- Test PR title 7 [\#2274](https://test.example.com/mrs/2274) ([test5-user](https://test.example.com/authors/test5-user))

## [v0.0.7](https://test.example.com/tags/v0.0.7) (14.03.2003)

Merged pull requests
--------------------
- Test PR title 6 [\#2264](https://test.example.com/mrs/2264) ([test-user](https://test.example.com/authors/test-user))

{{- if .TestingTag }}

## [testingtag](https://test.example.com/tags/testingtag) (13.03.2003)

{{- end }}

## [v0.0.6](https://test.example.com/tags/v0.0.6) (13.03.2003)

Merged pull requests
--------------------
- Test PR title 5 [\#2254](https://test.example.com/mrs/2254) ([test-user](https://test.example.com/authors/test-user))

## [v0.0.5](https://test.example.com/tags/v0.0.5) (12.03.2003)

Closed issues
-------------
- Test issue title 4 [\#1244](http://test.example.com/issues/1244)

## [v0.0.4](https://test.example.com/tags/v0.0.4) (11.03.2003)

Merged pull requests
--------------------
- Test PR title 3 [\#2234](https://test.example.com/mrs/2234) ([test-user](https://test.example.com/authors/test-user))

## [v0.0.3](https://test.example.com/tags/v0.0.3) (10.03.2003)

Closed issues
-------------
- Test issue title 2 [\#1227](http://test.example.com/issues/1227)

Merged pull requests
--------------------
- Test PR title 2 [\#2224](https://test.example.com/mrs/2224) ([test-user2](https://test.example.com/authors/test-user2))

## [v0.0.2](https://test.example.com/tags/v0.0.2) (09.03.2003)

Closed issues
-------------
- Test issue title 1 [\#1214](http://test.example.com/issues/1214)

Merged pull requests
--------------------
- Test PR title 1 [\#2214](https://test.example.com/mrs/2214) ([test-user](https://test.example.com/authors/test-user))

## [v0.0.1](https://test.example.com/tags/v0.0.1) (08.03.2003)

*This Changelog was automatically generated with [chagen unknown](https://github.com/artem-sidorenko/chagen)*
`

	input := struct {
		NewRelease     bool
		NewReleaseDate string
		TestingTag     bool
		SecondTag      bool
		ExcludedIssue  bool
	}{newRelease, time.Now().Format("02.01.2006"), testingTag, minorTags, excludedIssue}

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
		endpoint      string
	}

	tests := []struct {
		name                 string
		cliParams            cliParams
		repositoryExistsFail bool
		wantErr              error
		wantOutput           string
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
				filterExpr: `^v\d+\.0\.\d+$`,
			},
			wantOutput: genOutput(false, false, false, false),
		},
		{
			name: "With broken filter",
			cliParams: cliParams{
				filterExpr: "(abdc",
			},
			wantErr: errors.New("can't compile the regular expression: error parsing regexp: missing closing ): `(abdc`"), // nolint: lll
		},
		{
			name: "With customized labels",
			cliParams: cliParams{
				excludeLabels: "issue12, duplicate, question, invalid, wontfix, no changelog",
				newRelease:    "v10.10.0",
			},
			wantOutput: genOutput(true, false, true, true),
		},
		{
			name:                 "Repository not found",
			repositoryExistsFail: true,
			wantErr:              errors.New("project not found"),
		},
		{
			name: "With wrong endpoint type",
			cliParams: cliParams{
				endpoint: "wrongendpoint",
			},
			wantErr: errors.New("given endpoint isn't supported: wrongendpoint"),
		},
	}
	for _, tt := range tests {
		cliFlags := map[string]string{
			"file":     "-",
			"endpoint": "testconnector",
		}
		if tt.cliParams.endpoint != "" {
			cliFlags["endpoint"] = tt.cliParams.endpoint
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

		// avoid progress output
		progressOutput := &bytes.Buffer{}
		generate.ProgressWriter = progressOutput

		testconnector.RetTestingTag = true
		testconnector.RepositoryExistsFail = tt.repositoryExistsFail

		t.Run(tt.name, func(t *testing.T) {
			err := generate.Generate(ctx)
			out := output.String()

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				t.FailNow()
			}
			if out != tt.wantOutput {
				t.Errorf("Generate() output = %v, wantOutput %v", out, tt.wantOutput)
			}
		})
	}
}
