package api

import (
	"strings"
)

func splitDomainHost(hostname string) (domain, host string) {
	host = strings.Split(hostname, ".")[0]
	domain = strings.Replace(hostname, host+".", "", 1)
	return
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
