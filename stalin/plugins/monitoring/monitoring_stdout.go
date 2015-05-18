package monitoring

import (
	"net"

	. "stalin/plugins"
)

func init() {
	Plugin.Creator(newMonitoringStdout).Type("MonitoringStdout").Description("Run command in sandbox and resend stdout.").Register()
}

type MonitoringStdout struct {
	Log  *Logger  `json:"-"`
	Cmd  string   `json:"command" description:"Command for run"`
	Env  []string `json:"env" description:"Enviroment for 'cmd'"`
	Pwd  string   `json:"pwd" description:"Change dir before running 'cmd' (if empty sets to user home)"`
	User string   `json:"user" description:"Change gid and uid for process"`
}
