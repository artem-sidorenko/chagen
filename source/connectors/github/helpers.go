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

package github

import (
	"net/url"
	"path"

	"github.com/artem-sidorenko/chagen/source/connectors/helpers"
)

// formatErrorCode formats the error message for this connector
func formatErrorCode(query string, err error) error {
	return helpers.FormatErrorCode("GitHub", query, err)
}

// getTagURL returns the URL for a given tag.
// If alwaysUseReleaseURL is true: URL is provided for release page,
// even if it does not exist yet
func (c *Connector) getTagURL(tagName string, alwaysUseReleaseURL bool) (string, error) {
	release, resp, err := c.client.Repositories.GetReleaseByTag(c.context, c.Owner, c.Repo, tagName)
	if err != nil {
		// no release was found for this tag, this is no error for us
		if resp.StatusCode != 404 {
			return "", formatErrorCode("getTagURL", err)
		}
	}

	// if GitHub release for this tag was found -> use it
	// generate otherwise a link to the git tag view in the file tree
	var tagURL string
	if release != nil { // we got real release URL, use it
		tagURL = release.GetHTMLURL()
	} else { // build own URL
		u, err := url.Parse(c.ProjectURL)
		if err != nil {
			return "", err
		}

		if alwaysUseReleaseURL { // try to build own release url
			u.Path = path.Join(u.Path, "/releases/"+tagName)
		} else { // build tag url
			u.Path = path.Join(u.Path, "/tree/"+tagName)
		}
		tagURL = u.String()
	}
	return tagURL, nil
}

// GetNewTagURL returns the URL for a new tag, which does not exist yet
func (c *Connector) GetNewTagURL(TagName string) (string, error) {
	return c.getTagURL(TagName, c.NewTagUseReleaseURL)
}
