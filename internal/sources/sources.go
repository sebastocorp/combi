package sources

import (
	"fmt"
	"slices"
)

const (
	TypeRAW  = "RAW"
	TypeFILE = "FILE"
	TypeGIT  = "GIT"
	TypeK8S  = "K8S"
)

type SourceSetT struct {
	set map[string]SourceT
}

func (s *SourceSetT) Add(ops OptionsT) (err error) {
	switch ops.Type {
	case TypeRAW:
		{
			s.set[ops.Name], err = NewRawSource(ops)
		}
	case TypeFILE:
		{
			s.set[ops.Name], err = NewFileSource(ops)
		}
	case TypeGIT:
		{
			s.set[ops.Name], err = NewGitSource(ops)
		}
	case TypeK8S:
		{
			s.set[ops.Name], err = NewK8sSource(ops)
		}
	default:
		{
			err = fmt.Errorf("unsupported source type '%s'", ops.Type)
		}
	}
	return err
}

func (s *SourceSetT) Sync() (updated bool, err error) {
	var us []bool
	for sk := range s.set {
		var up bool
		up, err = s.set[sk].SyncConfig()
		if err != nil {
			return updated, err
		}
		us = append(us, up)
	}
	
	updated = slices.Contains(us, true)
	return updated, err
}

func (s *SourceSetT) Get(name string) ([]byte, error) {
	return nil, nil
}

type SourceT interface {
	GetName() string
	SyncConfig() (bool, error)
	GetConfig() ([]byte, error)
}

type OptionsT struct {
	Name string
	Type string
	Path string

	Raw  string
	File string
	Git  OptionsGitT
	K8s  OptionsK8sT
}

type OptionsGitT struct {
	SshKeyFilepath string
	Url            string
	Branch         string
	Filepath       string
}

type OptionsK8sT struct {
	InCluster      bool
	ConfigFilepath string
	MasterUrl      string
	Kind           string
	Namespace      string
	Name           string
	Key            string
}

func GetSource(ops OptionsT) (SourceT, error) {
	switch ops.Type {
	case TypeRAW:
		{
			return NewRawSource(ops)
		}
	case TypeFILE:
		{
			return NewFileSource(ops)
		}
	case TypeGIT:
		{
			return NewGitSource(ops)
		}
	case TypeK8S:
		{
			return NewK8sSource(ops)
		}
	}
	return nil, nil
}
