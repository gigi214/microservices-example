package service

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"
)

// BashExecService describes the service.
type BashExecService interface {
	// Add your methods here
	ExecCmd(ctx context.Context, cmd string) (stdOut string, stdErr string, exitCode int, err error)
}

type basicBashExecService struct{}

// NewBasicBashExecService returns a naive, stateless implementation of BashExecService.
func NewBasicBashExecService() BashExecService {
	return &basicBashExecService{}
}

// New returns a BashExecService with all of the expected middleware wired in.
func New(middleware []Middleware) BashExecService {
	var svc BashExecService = NewBasicBashExecService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}

var ErrInvalidCommand = errors.New("invalid command")

func (b *basicBashExecService) ExecCmd(ctx context.Context, cmd string) (stdOut string, stdErr string, exitCode int, err error) {
	stdOut = ""
	stdErr = ""
	exitCode = -999

	if len(cmd) < 3 {
		err = ErrInvalidCommand
		return
	}

	// Trim the command string to remove spaces at the beginning and ending
	cmd = strings.TrimSpace(cmd)

	// Split command by spaces in order to take the first token as the command, and the others as args
	splittedCommand := strings.Split(cmd, " ")

	var stdoutbb, stderrbb bytes.Buffer

	c := exec.Command(splittedCommand[0], splittedCommand[1:]...)
	c.Stdout = &stdoutbb
	c.Stderr = &stderrbb

	err = c.Run()

	exitCode = c.ProcessState.ExitCode()
	stdErr = stderrbb.String()
	stdOut = stdoutbb.String()

	// Return the output
	return stdOut, stdErr, exitCode, err
}
