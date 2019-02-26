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

	"github.com/artem-sidorenko/chagen/source/connectors/gitlab/internal/client"
	"github.com/artem-sidorenko/chagen/source/connectors/internal/testing/apitestdata"

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

func newProjectService() *ProjectsService {
	return &ProjectsService{
		ReturnValue: ReturnValue,
	}
}

func newTagsService() *TagsService {
	rtags := []*gitlab.Tag{}

	for _, tag := range apitestdata.Tags() {
		rtags = append(rtags, genTag(tag.Tag, tag.Commit, tag.Time))
	}

	return &TagsService{
		ReturnValue: ReturnValue,
		RetListTags: rtags,
	}
}

// New returns the configured simulated gitlab API client
func New(_ context.Context, _ string) *client.Client {
	return &client.Client{
		Projects: newProjectService(),
		Tags:     newTagsService(),
	}
}
