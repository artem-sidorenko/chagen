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

package client

import (
	gitlab "github.com/xanzy/go-gitlab"
)

// ProjectsService describes the methods we use from
// gitlab.ProjectsService
type ProjectsService interface {
	GetProject(
		pid interface{},
		options ...gitlab.OptionFunc,
	) (*gitlab.Project, *gitlab.Response, error)
}

// TagsService describes the methods we use from
// gitlab.TagsService
type TagsService interface {
	ListTags(
		pid interface{},
		opt *gitlab.ListTagsOptions,
		options ...gitlab.OptionFunc,
	) ([]*gitlab.Tag, *gitlab.Response, error)
}

// MergeRequestsService describes the methods we use from gitlab.MergeRequestsService
type MergeRequestsService interface {
	ListProjectMergeRequests(
		pid interface{},
		opt *gitlab.ListProjectMergeRequestsOptions,
		options ...gitlab.OptionFunc,
	) ([]*gitlab.MergeRequest, *gitlab.Response, error)
}

// CommitsService describes the methods we use from gitlab.CommitsService
type CommitsService interface {
	GetCommit(
		pid interface{},
		sha string,
		options ...gitlab.OptionFunc,
	) (*gitlab.Commit, *gitlab.Response, error)
}

// IssuesService describes the methods we use from gitlab.IssuesService
type IssuesService interface {
	ListProjectIssues(
		pid interface{},
		opt *gitlab.ListProjectIssuesOptions,
		options ...gitlab.OptionFunc,
	) ([]*gitlab.Issue, *gitlab.Response, error)
}

// Client wraps the gitlab.Client with interfaces we are using
type Client struct {
	Projects      ProjectsService
	Tags          TagsService
	MergeRequests MergeRequestsService
	Commits       CommitsService
	Issues        IssuesService
}
