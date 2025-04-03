package sources

import (
	"fmt"
	"slices"
)

const (
	TypeFILERAW = "FILERAW"
	TypeFILE    = "FILE"
	TypeGIT     = "GIT"
	TypeK8S     = "K8S"
)

type SetT struct {
	size int
	ss   []SourceT
}

type SourceT interface {
	Name() string
	sync() (bool, error)
	get() ([]byte, error)
}

type OptionsT struct {
	Name string
	Type string
	Path string
	Cred any

	Raw  string
	File string
	Git  OptionsGitT
	K8s  OptionsK8sT
}

func NewSet() (s *SetT, err error) {
	s = &SetT{}
	return s, err
}

func (s *SetT) Add(ops OptionsT) (err error) {
	switch ops.Type {
	case TypeFILE, TypeFILERAW:
		{
			var src *FileSourceT
			src, err = NewFileSource(ops)
			if err != nil {
				return err
			}
			s.ss = append(s.ss, src)
		}
	case TypeGIT:
		{
			var src *GitSourceT
			src, err = NewGitSource(ops)
			if err != nil {
				return err
			}
			s.ss = append(s.ss, src)
		}
	case TypeK8S:
		{
			var src *K8sSourceT
			src, err = NewK8sSource(ops)
			if err != nil {
				return err
			}
			s.ss = append(s.ss, src)
		}
	default:
		{
			err = fmt.Errorf("unsupported source type '%s'", ops.Type)
			return err
		}
	}
	s.size++

	return err
}

func (s *SetT) Length() int {
	return s.size
}

func (s *SetT) Sync() (updated bool, err error) {
	var us []bool
	for sk := range s.ss {
		var up bool
		up, err = s.ss[sk].sync()
		if err != nil {
			return updated, err
		}
		us = append(us, up)
	}

	updated = slices.Contains(us, true)
	return updated, err
}

func (s *SetT) Get(name string) ([]byte, error) {
	return nil, nil
}

func (s *SetT) GetByIndex(index int) ([]byte, error) {
	if index >= s.size || index < 0 {
		return nil, fmt.Errorf("index out of bounds")
	}

	return s.ss[index].get()
}

func GetSource(ops OptionsT) (SourceT, error) {
	switch ops.Type {
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
