package git

import (
	"os"
	"path/filepath"
	"reflect"

	"combi/api/v1alpha3"
	"combi/internal/config"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type GitSourceT struct {
	name       string
	srcConfig  string
	storConfig string

	repo repoT
}

type repoT struct {
	syncPath string
	sshKey   string
	url      string
	branch   string
}

func NewGitSource(srcConf v1alpha3.SourceConfigT, srcpath string) (s *GitSourceT, err error) {
	s = &GitSourceT{
		name:       srcConf.Name,
		srcConfig:  filepath.Join(srcpath, "sync/repo", srcConf.Git.Filepath),
		storConfig: filepath.Join(srcpath, filepath.Base(srcConf.Git.Filepath)),

		repo: repoT{
			syncPath: filepath.Join(srcpath, "sync/repo"),
			sshKey:   srcConf.Git.SshKeyFilepath,
			url:      srcConf.Git.SshUrl,
			branch:   srcConf.Git.Branch,
		},
	}

	if _, err = os.Stat(s.repo.sshKey); err != nil {
		return s, err
	}

	err = os.MkdirAll(filepath.Join(srcpath, "sync"), 0777)
	if err != nil {
		return s, err
	}

	return s, err
}

func (s *GitSourceT) GetName() string {
	return s.name
}

func (s *GitSourceT) SyncConfig() (updated bool, err error) {
	if _, err = os.Stat(s.srcConfig); !os.IsNotExist(err) {
		if err = os.RemoveAll(s.repo.syncPath); err != nil {
			return updated, err
		}
	}

	publicSshKey, err := ssh.NewPublicKeysFromFile("git", s.repo.sshKey, "")
	if err != nil {
		return updated, err
	}

	_, err = git.PlainClone(s.repo.syncPath, false, &git.CloneOptions{
		URL:           s.repo.url,
		Depth:         1,
		ReferenceName: plumbing.NewBranchReferenceName(s.repo.branch),
		SingleBranch:  true,
		Auth:          publicSshKey,
	})
	if err != nil {
		return updated, err
	}

	srcBytes, err := os.ReadFile(s.srcConfig)
	if err != nil {
		return updated, err
	}

	storBytes, err := os.ReadFile(s.storConfig)
	if err != nil {
		if os.IsNotExist(err) {
			updated = true
			err = os.WriteFile(s.storConfig, srcBytes, 0777)
			if err != nil {
				return updated, err
			}
		}
		return updated, err
	}

	if !reflect.DeepEqual(srcBytes, storBytes) {
		updated = true
		err = os.WriteFile(s.storConfig, srcBytes, 0777)
		if err != nil {
			return updated, err
		}
	}

	return updated, err
}

func (s *GitSourceT) GetConfig() (conf []byte, err error) {
	if conf, err = os.ReadFile(s.storConfig); err != nil {
		return conf, err
	}

	conf = config.ExpandEnv(conf)

	return conf, err
}
