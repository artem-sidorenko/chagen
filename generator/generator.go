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

// Package generator provides the capabilities for changelog generation
package generator

import (
	"io"
	"text/template"

	"github.com/artem-sidorenko/chagen/data"
)

const changelogTemplate = `Changelog
=========
{{ range .Releases}}
## [{{.Release}}]({{.ReleaseURL}}) ({{.Date}})

Closed issues
-------------
{{- range .Issues}}
- {{.Name}} [\#{{.ID}}]({{.URL}})
{{- end}}
{{ end}}`

// Generator is resposible for generation of Changelogs.
// Each data field represents the data structure, which is consumed by the template.
type Generator struct {
	Releases []data.Release
}

// Render the content via template and write it to wr.
// It returns the result of template complication
func (g *Generator) Render(wr io.Writer) error {
	t := template.Must(template.New("Changelog template").Parse(changelogTemplate))
	return t.Execute(wr, g)
}
