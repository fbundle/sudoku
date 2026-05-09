package http_transport

import "strings"

func join(prefix, path string) string {
	return strings.Trim(prefix+"/"+strings.Trim(path, "/"), "/")
}
