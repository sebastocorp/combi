package sources

import (
	"fmt"
	"slices"
)

const (
	TypeFILE = "FILE"
	TypeGIT  = "GIT"
	TypeK8S  = "K8S"
)

type SetT struct {
	size int
	ss   []SourceT
}

type SourceT interface {
	getName() string
	getData() (SourceDataT, error)
	sync() (bool, error)
}

type OptionsT struct {
	Name    string
	SrcType string
	EncType string
	WorkDir string
	CredRef any

	File string
	Git  OptionsGitT
	K8s  OptionsK8sT
}

type SourceDataT struct {
	Name    string
	SrcType string
	EncType string
	Data    []byte
}

func NewSet() (s *SetT, err error) {
	s = &SetT{}
	return s, err
}

func (s *SetT) Add(ops OptionsT) (err error) {
	switch ops.SrcType {
	case TypeFILE:
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
			err = fmt.Errorf("unsupported source type '%s'", ops.SrcType)
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

func (s *SetT) GetByName(name string) (SourceDataT, error) {
	for index := range s.ss {
		if s.ss[index].getName() == name {
			return s.ss[index].getData()
		}
	}

	return SourceDataT{}, fmt.Errorf("source '%s' not found in set", name)
}

func (s *SetT) GetByIndex(index int) (SourceDataT, error) {
	if index >= s.size || index < 0 {
		return SourceDataT{}, fmt.Errorf("index out of bounds")
	}

	return s.ss[index].getData()
}
