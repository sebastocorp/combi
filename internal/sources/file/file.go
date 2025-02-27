package file

import (
	"combi/api/v1alpha4"
	"combi/internal/utils"
	"os"
	"path/filepath"
	"reflect"
)

type FileSourceT struct {
	name       string
	srcConfig  string
	storConfig string
}

func NewFileSource(srcConf v1alpha4.SourceConfigT, srcpath string) (s *FileSourceT, err error) {
	s = &FileSourceT{
		name:       srcConf.Name,
		srcConfig:  srcConf.File,
		storConfig: filepath.Join(srcpath, filepath.Base(srcConf.File)),
	}

	return s, err
}

func (s *FileSourceT) GetName() string {
	return s.name
}

func (s *FileSourceT) SyncConfig() (updated bool, err error) {
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

func (s *FileSourceT) GetConfig() (conf []byte, err error) {
	if conf, err = os.ReadFile(s.storConfig); err != nil {
		return conf, err
	}

	conf = utils.ExpandEnv(conf)

	return conf, err
}
