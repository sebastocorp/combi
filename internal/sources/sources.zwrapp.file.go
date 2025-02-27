package sources

import (
	"combi/internal/utils"
	"os"
	"path/filepath"
	"reflect"
)

type FileSourceT struct {
	name    string
	tmpPath string
	file    string
}

func NewFileSource(ops OptionsT) (s *FileSourceT, err error) {
	s = &FileSourceT{
		name:    ops.Name,
		tmpPath: ops.Path,
		file:    ops.File,
	}

	return s, err
}

func (s *FileSourceT) GetName() string {
	return s.name
}

func (s *FileSourceT) SyncConfig() (updated bool, err error) {
	srcBytes, err := os.ReadFile(s.file)
	if err != nil {
		return updated, err
	}

	storConfig := filepath.Join(s.tmpPath, filepath.Base(s.file))
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

func (s *FileSourceT) GetConfig() (conf []byte, err error) {
	storConfig := filepath.Join(s.tmpPath, filepath.Base(s.file))
	if conf, err = os.ReadFile(storConfig); err != nil {
		return conf, err
	}

	conf = utils.ExpandEnv(conf)

	return conf, err
}
