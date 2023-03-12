package rest

import (
	"regexp"
)

func getPathID(path string) (string, bool) {
	re := regexp.MustCompile("(?:/)([[:digit:]]+)$")

	if match := re.FindStringSubmatch(path); len(match) == 2 {
		return match[1], true
	}

	return "", false
}
