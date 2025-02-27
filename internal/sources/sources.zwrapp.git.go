package sources

import (
	"combi/internal/utils"
	"os"
	"path/filepath"
	"reflect"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type GitSourceT struct {
	name    string
	tmpPath string

	repo repoT
}

type repoT struct {
	sshKey string
	url    string
	branch string
	file   string
}

func NewGitSource(ops OptionsT) (s *GitSourceT, err error) {
	s = &GitSourceT{
		name:    ops.Name,
		tmpPath: ops.Path,

		repo: repoT{
			sshKey: ops.Git.SshKeyFilepath,
			url:    ops.Git.Url,
			branch: ops.Git.Branch,
			file:   ops.Git.Filepath,
		},
	}

	if _, err = os.Stat(s.repo.sshKey); err != nil {
		return s, err
	}

	err = os.MkdirAll(filepath.Join(ops.Path, "sync"), 0777)
	if err != nil {
		return s, err
	}

	return s, err
}

func (s *GitSourceT) GetName() string {
	return s.name
}

func (s *GitSourceT) SyncConfig() (updated bool, err error) {
	syncPath := filepath.Join(s.tmpPath, "sync/repo")
	srcConfig := filepath.Join(syncPath, s.repo.file)
	if _, err = os.Stat(srcConfig); !os.IsNotExist(err) {
		if err = os.RemoveAll(syncPath); err != nil {
			return updated, err
		}
	}

	publicSshKey, err := ssh.NewPublicKeysFromFile("git", s.repo.sshKey, "")
	if err != nil {
		return updated, err
	}

	_, err = git.PlainClone(syncPath, false, &git.CloneOptions{
		URL:           s.repo.url,
		Depth:         1,
		ReferenceName: plumbing.NewBranchReferenceName(s.repo.branch),
		SingleBranch:  true,
		Auth:          publicSshKey,
	})
	if err != nil {
		return updated, err
	}

	srcBytes, err := os.ReadFile(srcConfig)
	if err != nil {
		return updated, err
	}

	storConfig := filepath.Join(s.tmpPath, filepath.Base(s.repo.file))
	storBytes, err := os.ReadFile(storConfig)
	if err != nil {
		if os.IsNotExist(err) {
			updated = true
			err = os.WriteFile(storConfig, srcBytes, 0777)
			if err != nil {
				return updated, err
			}
		}
		return updated, err
	}

	if !reflect.DeepEqual(srcBytes, storBytes) {
		updated = true
		err = os.WriteFile(storConfig, srcBytes, 0777)
		if err != nil {
			return updated, err
		}
	}

	return updated, err
}

func (s *GitSourceT) GetConfig() (conf []byte, err error) {
	storConfig := filepath.Join(s.tmpPath, filepath.Base(s.repo.file))
	if conf, err = os.ReadFile(storConfig); err != nil {
		return conf, err
	}

	conf = utils.ExpandEnv(conf)

	return conf, err
}
