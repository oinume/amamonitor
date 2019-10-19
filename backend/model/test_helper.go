package model

import (
	"strings"
)

func ReplaceToTestDBURL(dbURL string) string {
	if strings.Contains(dbURL, "/amamonitor?") {
		return strings.Replace(dbURL, "/amamonitor?", "/amamonitor_test?", 1)
	}
	return dbURL
}
