package combi

import (
	"os"
	"os/exec"

	"combi/api/v1alpha3"
)



type ActionT struct {
	Name string `json:"name"`
	On   string `json:"on"`

	cmd []string `json:"-"`
}

func NewAction(action v1alpha3.ActionConfigT) ActionT {
	return ActionT{
		Name: action.Name,
		On:   action.On,

		cmd: action.Command,
	}
}

func (a *ActionT) Exec() error {
	command := exec.Command(a.cmd[0], a.cmd[1:]...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}
