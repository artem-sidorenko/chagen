package gitlab

import (
	"net/url"
	"path"

	"github.com/artem-sidorenko/chagen/source/connectors/helpers"
)

// formatErrorCode formats the error message for this connector
func formatErrorCode(query string, err error) error { // nolint: unparam
	return helpers.FormatErrorCode("GitLab", query, err)
}

// getTagURL returns the URL for a given tag
func (c *Connector) getTagURL(tagName string) (string, error) {
	u, err := url.Parse(c.ProjectURL)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, "/tags/"+tagName)
	return u.String(), nil
}
