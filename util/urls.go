package util

import "strings"

func GetProtocol(url string) string {
	if strings.HasPrefix(url, "ssh://") || strings.HasPrefix(url, "git@") {
		return "ssh"
	}
	return "http"
}
