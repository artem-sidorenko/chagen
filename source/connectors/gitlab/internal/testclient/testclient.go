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

package testclient

import (
	"context"
	"fmt"
	"time"

	"github.com/artem-sidorenko/chagen/source/connectors/gitlab/internal/client"

	gitlab "github.com/xanzy/go-gitlab"
)

// ReturnValueStr represents the possible error return values of API
// if some field is true - error is return, otherise not
type ReturnValueStr struct {
	ProjectsServiceGetProjectRespCode int
	ProjectsServiceGetProjectErr      bool
	TagsServiceListTagsErr            bool
}

// ReturnValue controls the error return values of API for testclient instances
var ReturnValue = ReturnValueStr{} // nolint: gochecknoglobals

// ProjectsService simulates the gitlab.ProjectsService
type ProjectsService struct {
	ReturnValue ReturnValueStr
}

// GetProject simulates the (gitlab.ProjectsService).GetProject call
func (p *ProjectsService) GetProject(
	_ interface{},
	_ ...gitlab.OptionFunc,
) (*gitlab.Project, *gitlab.Response, error) {

	respCode := 200
	if p.ReturnValue.ProjectsServiceGetProjectRespCode != 0 {
		respCode = p.ReturnValue.ProjectsServiceGetProjectRespCode
	}

	response := genResponse(respCode)

	if p.ReturnValue.ProjectsServiceGetProjectErr {
		return nil, response, fmt.Errorf("can't fetch the repo data")
	}

	return nil, response, nil
}

// TagsService sumulates the gitlab.TagsService
type TagsService struct {
	RetListTags []*gitlab.Tag
	ReturnValue ReturnValueStr
}

// ListTags simulates the (gitlab.TagsService).ListTags call
func (t *TagsService) ListTags(
	_ interface{},
	opt *gitlab.ListTagsOptions,
	_ ...gitlab.OptionFunc,
) ([]*gitlab.Tag, *gitlab.Response, error) {

	if t.ReturnValue.TagsServiceListTagsErr {
		return nil, nil, fmt.Errorf("can't fetch the tags")
	}

	resp, start, end := calcPaging(opt.Page, opt.PerPage, len(t.RetListTags))

	return t.RetListTags[start:end], resp, nil
}

type gitLabTagsServiceInput struct {
	tag            string
	commit         string
	time           time.Time
	releasePresent bool
}

// New returns the configured simulated gitlab API client
func New(_ context.Context, _ string) *client.Client {
	p := &ProjectsService{
		ReturnValue: ReturnValue,
	}

	t := &TagsService{
		ReturnValue: ReturnValue,
	}
	ts := []gitLabTagsServiceInput{
		{"v0.0.1", "7d84cdb2f7c2d4619cda4b8adeb1897097b5c8fc", time.Unix(2047083647, 0), true},
		{"v0.0.2", "b3622b516b8ad70ce5dc3fa422fb90c3b58fa9da", time.Unix(2047183647, 0), false},
		{"v0.0.3", "52f214dc3bf6c0e2a87eae6eab363a317c5a665f", time.Unix(2047283647, 0), true},
		{"v0.0.4", "d4ff341587bc80a9c897c28340df9fe8f9fc6309", time.Unix(2047383647, 0), false},
		{"v0.0.5", "746e45ea014e257bcb7caa2c100ed1e5f63ed234", time.Unix(2047483647, 0), false},
		{"v0.0.6", "ddde800c451bae606713ae0f8418badcf31db120", time.Unix(2047583647, 0), false},
		{"v0.0.7", "d21438494dd0722c1d13dc496ae1f60fb85084c1", time.Unix(2047683647, 0), true},
		{"v0.0.8", "8d8d817a530bc1c3f792d9508c187b5769c434c5", time.Unix(2047783647, 0), false},
		{"v0.0.9", "fc9f16ecc043e3fe422834cd127311d11d423668", time.Unix(2047883647, 0), false},
		{"v0.1.0", "dbbf36ffaae700a2ce03ef849d6f944031f34b95", time.Unix(2047983647, 0), true},
		{"v0.1.1", "fc5d68ff1cf691e09f6ead044813274953c9b843", time.Unix(2048083647, 0), true},
		{"v0.1.2", "d8351413f688c96c2c5d6fe58ebf5ac17f545bc0", time.Unix(2048183647, 0), true},
	}
	for _, tag := range ts {
		t.RetListTags = append(t.RetListTags, genTag(tag.tag, tag.commit, tag.time))
	}

	return &client.Client{
		Projects: p,
		Tags:     t,
	}
}