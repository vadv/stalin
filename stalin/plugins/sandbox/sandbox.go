package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"

	. "stalin/plugins"
	. "stalin/supervisor"
)

func init() {
	Plugin.Creator(newSandbox).Type("Sandbox").Description("Run command in sandbox").Register()
}

type Sandbox struct {
	Log  *Logger  `json:"-"`
	Cmd  string   `json:"cmd" description:"Path to command with args"`
	Env  []string `json:"env" description:"Enviroment for 'cmd'"`
	Dir  string   `json:"pwd" description:"Change dir before running 'cmd' (if empty sets to user home)"`
	User string   `json:"user" description:"Change gid and uid for process"`
}

func newSandbox(name string) PluginInterface {
	return &Sandbox{
		Log: NewLog(name),
		Env: []string{fmt.Sprintf("PATH=%s", os.Getenv("PATH"))},
	}
}

func (s *Sandbox) cmd() (*exec.Cmd, error) {
	cmd := strings.Split(s.Cmd, " ")
	if len(cmd) == 0 {
		return nil, fmt.Errorf("Command '%s' is empty", s.Cmd)
	}
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Dir = s.Dir
	command.Env = s.Env
	if s.User != "" {
		if sandboxUser, err := user.Lookup(s.User); err != nil {
			return nil, err
		} else {
			uid, err := strconv.ParseUint(sandboxUser.Uid, 10, 32)
			if err != nil {
				return nil, err
			}
			gid, err := strconv.ParseUint(sandboxUser.Gid, 10, 32)
			if err != nil {
				return nil, err
			}
			command.SysProcAttr = &syscall.SysProcAttr{}
			command.SysProcAttr.Credential = &syscall.Credential{
				Uid: uint32(uid),
				Gid: uint32(gid),
			}
			if command.Dir == "" {
				command.Dir = sandboxUser.HomeDir
			}
		}
	}
	return command, nil
}

func (s *Sandbox) Start() error {
	cmd, err := s.cmd()
	if err != nil {
		return err
	}
	supervisor := NewSupervisor(cmd)
	stdout := make(chan string, 10)
	stderr := make(chan string, 10)
	errchan := make(chan error, 10)
	go supervisor.Run(errchan, stdout, stderr)
	for {
		select {
		case err := <-errchan:
			if err != nil {
				s.Log.Error("[fatal]: '%v'", err)
			}
		case out := <-stdout:
			s.Log.Info("[stdout]: %s", out)
		case out := <-stderr:
			s.Log.Error("[stderr]: %s", out)
		}
	}
}
