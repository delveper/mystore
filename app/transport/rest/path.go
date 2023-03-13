package rest

import (
	"path"
	"strconv"
)

func getPathID(pth string) (id int, ok bool) {
	base := path.Base(path.Clean(pth))

	id, err := strconv.Atoi(base)
	if err != nil {
		return 0, false
	}

	return id, true
}
