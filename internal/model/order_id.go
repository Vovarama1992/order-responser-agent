package model

import (
	"regexp"
)

var orderIDRegexp = regexp.MustCompile(`/projects/(\d+)/`)

func ExtractOrderID(url string) string {
	match := orderIDRegexp.FindStringSubmatch(url)

	if len(match) != 2 {
		return ""
	}

	return match[1]
}
