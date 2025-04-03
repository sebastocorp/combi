package credentials

import (
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type SshKeyT struct {
	PublicKey *ssh.PublicKeys
}

type OptionsSshKeyT struct {
	User     string
	SshKey   string
	Password string
}

func NewSshKey(ops OptionsSshKeyT) (sc *SshKeyT, err error) {
	sc = &SshKeyT{}
	if ops.User == "" {
		ops.User = "git"
	}

	sc.PublicKey, err = ssh.NewPublicKeysFromFile(ops.User, ops.SshKey, ops.Password)

	return sc, err
}
