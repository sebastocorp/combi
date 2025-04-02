package credentials

import "fmt"

const (
	TypeK8S    = "K8S"
	TypeSSHKEY = "SSK_KEY"
)

type CredentialSetT struct {
	set map[string]any
}

type OptionsT struct {
	Name   string
	Type   string
	Kube   OptionsKubeT
	SshKey OptionsSshKeyT
}

func (cs *CredentialSetT) Add(ops OptionsT) (err error) {
	switch ops.Type {
	case TypeK8S:
		{
			cs.set[ops.Name], err = NewKube(ops.Kube)
		}
	case TypeSSHKEY:
		{
			cs.set[ops.Name], err = NewSshKey(ops.SshKey)
		}
	default:
		{
			err = fmt.Errorf("unsupported credential type '%s'", ops.Type)
		}
	}
	return err
}

func (cs *CredentialSetT) Get(name string) any {
	return cs.set[name]
}
