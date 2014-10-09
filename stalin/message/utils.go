package message

import (
	"os"
)

func CurrentHost() string {
	host, _ := os.Hostname()
	return host
}
