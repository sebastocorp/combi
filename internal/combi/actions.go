package combi

import (
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

func (a *ActionT) Exec() (out []byte, err error) {
	return exec.Command(a.cmd[0], a.cmd[1:]...).CombinedOutput()
}
