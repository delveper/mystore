package rest

import (
	"path"
	"regexp"
	"strings"
)

func getPathID(path string) (string, bool) {
	re := regexp.MustCompile("(?:/)([[:digit:]]+)$")

	if match := re.FindStringSubmatch(path); len(match) == 2 {
		return match[1], true
	}

	return "", false
}

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
