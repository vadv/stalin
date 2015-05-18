package supervisor

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"syscall"
	"time"
)

const timeSleepBeforeStart = 1

type supervisor struct {
	inCmd   *exec.Cmd
	workCmd *exec.Cmd
	status  *syscall.WaitStatus
	errChan chan error
	stdOut  chan string
	stdErr  chan string
}

func NewSupervisor(cmd *exec.Cmd) *supervisor {
	return &supervisor{
		inCmd:   cmd,
		errChan: make(chan error),
		stdOut:  make(chan string),
		stdErr:  make(chan string),
	}
}

func (s *supervisor) Run(outErrChan chan error, stdOut, stdErr chan string) error {
	go s.run()
	for {
		select {
		case out := <-s.stdOut:
			stdOut <- out
		case out := <-s.stdErr:
			stdErr <- out
		case err := <-s.errChan:
			if err != nil {
				s.setWaitStatus(err)
				outErrChan <- err
			} else {
				outErrChan <- fmt.Errorf("exit status 0")
			}
			time.Sleep(time.Duration(timeSleepBeforeStart) * time.Second)
			go s.run()
		}
	}
	return nil
}

func (s *supervisor) run() {
	s.kill()
	s.make()
	s.feedOut()
	if err := s.workCmd.Start(); err != nil {
		s.errChan <- err
		return
	}
	s.errChan <- s.workCmd.Wait()
}

func (s *supervisor) make() {
	s.workCmd = &exec.Cmd{
		Path:        s.inCmd.Path,
		Dir:         s.inCmd.Dir,
		Args:        s.inCmd.Args,
		Env:         s.inCmd.Env,
		SysProcAttr: s.inCmd.SysProcAttr,
	}
}

func (s *supervisor) kill() {
	if s.workCmd != nil && s.workCmd.Process != nil {
		if s.status != nil && !s.status.Exited() {
			s.errChan <- s.workCmd.Process.Kill()
		}
	}
}

func (s *supervisor) setWaitStatus(err error) {
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			s.status = &status
		}
	}
}

func (s *supervisor) feedOut() {
	stdout, _ := s.workCmd.StdoutPipe()
	go feed(stdout, s.stdOut)
	stderr, _ := s.workCmd.StderrPipe()
	go feed(stderr, s.stdErr)
}

func feed(pipe io.Reader, sink chan<- string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		sink <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return
	}
}
