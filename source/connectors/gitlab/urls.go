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

package gitlab

import (
	"net/url"
	"path"
)

// getTagURL returns the URL for a given tag
func (c *Connector) getTagURL(tagName string) (string, error) {
	u, err := url.Parse(c.ProjectURL)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, "/tags/"+tagName)
	return u.String(), nil
}

// getUsernameURL returns the URL for a given username
func (c *Connector) getUsernameURL(username string) (string, error) {
	u, err := url.Parse(c.ProjectURL)
	if err != nil {
		return "", err
	}

	u.Path = path.Join("/" + username)
	return u.String(), nil
}

// GetNewTagURL returns the URL for a new tag, which does not exist yet
func (c *Connector) GetNewTagURL(TagName string) (string, error) {
	return c.getTagURL(TagName)
}
