package gitlab

import (
	"github.com/artem-sidorenko/chagen/datasource/connectors/helpers"
)

// formatErrorCode formats the error message for this connector
func formatErrorCode(query string, err error) error { // nolint: unparam
	return helpers.FormatErrorCode("GitLab", query, err)
}
