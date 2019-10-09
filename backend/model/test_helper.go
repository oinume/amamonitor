package model

import (
	"strings"
	"testing"
)

func ReplaceToTestDBURL(t *testing.T, dbURL string) string {
	if strings.Contains(dbURL, "/amamonitor?") {
		return strings.Replace(dbURL, "/amamonitor?", "/amamonitor_test?", 1)
	}
	return dbURL
}
