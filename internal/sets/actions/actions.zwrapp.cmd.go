package actions

import (
	"bytes"
	"os/exec"
)

type CmdActionT struct {
	on string

	cmd []string
}

type OptionsCmdT struct {
	Cmd []string
}

func newCmdAction(ops OptionsT) (a *CmdActionT, err error) {
	a = &CmdActionT{
		on:  ops.On,
		cmd: ops.Cmd.Cmd,
	}
	return a, err
}

func (a *CmdActionT) getOn() string {
	return a.on
}

func (a *CmdActionT) exec() (r ActionResultT, err error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(a.cmd[0], a.cmd[1:]...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	r.Stdout = stdout.String()
	r.Stderr = stderr.String()
	return r, err
}
