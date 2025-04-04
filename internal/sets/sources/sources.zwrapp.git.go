package sources

import (
	"combi/internal/sets/credentials"
	"combi/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitSourceT struct {
	name    string
	encType string
	workDir string
	credRef *credentials.SshKeyT

	repoURL    string
	repoBranch string
	repoFile   string
}

type OptionsGitT struct {
	Url      string
	Branch   string
	Filepath string
}

func NewGitSource(ops OptionsT) (s *GitSourceT, err error) {
	s = &GitSourceT{
		name:    ops.Name,
		encType: ops.EncType,
		workDir: ops.WorkDir,
		credRef: ops.CredRef.(*credentials.SshKeyT),

		repoURL:    ops.Git.Url,
		repoBranch: ops.Git.Branch,
		repoFile:   ops.Git.Filepath,
	}

	switch ops.CredRef.(type) {
	case *credentials.SshKeyT:
		s.credRef = ops.CredRef.(*credentials.SshKeyT)
	default:
		err = fmt.Errorf("wrong credential type in '%s' source, must be SSH_KEY", ops.Name)
		return s, err
	}

	err = os.MkdirAll(filepath.Join(s.workDir, "sync"), 0644)
	if err != nil {
		return s, err
	}

	return s, err
}

func (s *GitSourceT) getName() string {
	return s.name
}

func (s *GitSourceT) getData() (srcd SourceDataT, err error) {
	srcd.Name = s.name
	srcd.SrcType = TypeGIT
	srcd.EncType = s.encType

	storConfig := filepath.Join(s.workDir, filepath.Base(s.repoFile))
	if srcd.Data, err = os.ReadFile(storConfig); err != nil {
		return srcd, err
	}
	srcd.Data = utils.ExpandEnv(srcd.Data)

	return srcd, err
}

func (s *GitSourceT) sync() (updated bool, err error) {
	syncPath := filepath.Join(s.workDir, "sync", "repo")
	srcConfig := filepath.Join(syncPath, s.repoFile)
	if _, err = os.Stat(srcConfig); !os.IsNotExist(err) {
		if err = os.RemoveAll(syncPath); err != nil {
			return updated, err
		}
	}

	_, err = git.PlainClone(syncPath, false, &git.CloneOptions{
		URL:           s.repoURL,
		Depth:         1,
		ReferenceName: plumbing.NewBranchReferenceName(s.repoBranch),
		SingleBranch:  true,
		Auth:          s.credRef.PublicKey,
	})
	if err != nil {
		return updated, err
	}

	srcBytes, err := os.ReadFile(srcConfig)
	if err != nil {
		return updated, err
	}

	storConfig := filepath.Join(s.workDir, filepath.Base(s.repoFile))
	storBytes, err := os.ReadFile(storConfig)
	if err != nil {
		if os.IsNotExist(err) {
			updated = true
			err = os.WriteFile(storConfig, srcBytes, 0755)
			if err != nil {
				return updated, err
			}
		}
		return updated, err
	}

	if !reflect.DeepEqual(srcBytes, storBytes) {
		updated = true
		err = os.WriteFile(storConfig, srcBytes, 0755)
		if err != nil {
			return updated, err
		}
	}

	return updated, err
}
