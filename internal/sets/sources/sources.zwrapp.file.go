package sources

import (
	"combi/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type FileSourceT struct {
	name    string
	encType string
	workDir string

	srcType string
	file    string
}

func NewFileSource(ops OptionsT) (s *FileSourceT, err error) {
	s = &FileSourceT{
		name:    ops.Name,
		encType: ops.EncType,
		workDir: ops.WorkDir,

		srcType: ops.SrcType,
	}

	err = os.MkdirAll(s.workDir, 0644)
	if err != nil {
		return s, err
	}

	switch ops.SrcType {
	case TypeFILE:
		{

			if !utils.FileExists(ops.File) {
				err = fmt.Errorf("file '%s' does not exist", ops.File)
				return s, err
			}
			s.file = ops.File
		}
	case TypeFILERAW:
		{
			srcPath := filepath.Join(s.workDir, "sync")
			err = os.MkdirAll(srcPath, 0644)
			if err != nil {
				return s, err
			}

			s.file = filepath.Join(srcPath, strings.Join([]string{"fileraw", strings.ToLower(s.encType), "txt"}, "."))
			if err = os.WriteFile(s.file, []byte(ops.File), 0644); err != nil {
				return s, err
			}
		}
	default:
		{
			err = fmt.Errorf("unsupported source type '%s'", ops.SrcType)
		}
	}

	return s, err
}

func (s *FileSourceT) getName() string {
	return s.name
}

func (s *FileSourceT) getData() (srcd SourceDataT, err error) {
	srcd.Name = s.name
	srcd.SrcType = s.srcType
	srcd.EncType = s.encType

	storConfig := filepath.Join(s.workDir, filepath.Base(s.file))
	if srcd.Data, err = os.ReadFile(storConfig); err != nil {
		return srcd, err
	}
	srcd.Data = utils.ExpandEnv(srcd.Data)

	return srcd, err
}

func (s *FileSourceT) sync() (updated bool, err error) {
	srcBytes, err := os.ReadFile(s.file)
	if err != nil {
		return updated, err
	}

	storConfig := filepath.Join(s.workDir, filepath.Base(s.file))
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

	if !reflect.DeepEqual(storBytes, srcBytes) {
		updated = true
		err = os.WriteFile(storConfig, srcBytes, 0755)
		if err != nil {
			return updated, err
		}
	}

	return updated, err
}
