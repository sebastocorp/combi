package credentials

import "fmt"

const (
	TypeK8S    = "K8S"
	TypeSSHKEY = "SSH_KEY"
)

type SetT struct {
	cs map[string]any
}

type OptionsT struct {
	Name   string
	Type   string
	Kube   OptionsKubeT
	SshKey OptionsSshKeyT
}

func NewSet() (set *SetT, err error) {
	set = &SetT{
		cs: make(map[string]any),
	}
	return set, err
}

func (s *SetT) Add(ops OptionsT) (err error) {
	switch ops.Type {
	case TypeK8S:
		{
			s.cs[ops.Name], err = NewKube(ops.Kube)
		}
	case TypeSSHKEY:
		{
			s.cs[ops.Name], err = NewSshKey(ops.SshKey)
		}
	default:
		{
			err = fmt.Errorf("unsupported credential type '%s'", ops.Type)
		}
	}
	return err
}

func (s *SetT) Get(name string) any {
	return s.cs[name]
}
