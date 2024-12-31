package file

import (
	"combi/api/v1alpha3"
	"combi/internal/config"
	"os"
	"path/filepath"
	"reflect"
)

type FileSourceT struct {
	srcConfig  string
	storConfig string
}

func NewFileSource(srcConf v1alpha3.SourceConfigT, srcpath string) (s *FileSourceT, err error) {
	s = &FileSourceT{
		srcConfig:  srcConf.File,
		storConfig: filepath.Join(srcpath, filepath.Base(srcConf.File)),
	}

	var configBytes []byte
	configBytes, err = os.ReadFile(s.srcConfig)
	if err != nil {
		return s, err
	}

	err = os.WriteFile(s.storConfig, configBytes, 0777)
	if err != nil {
		return s, err
	}

	return s, err
}

func (s *FileSourceT) SyncConfig() (updated bool, err error) {
	srcBytes, err := os.ReadFile(s.srcConfig)
	if err != nil {
		return updated, err
	}

	storBytes, err := os.ReadFile(s.storConfig)
	if err != nil {
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

	conf = config.ExpandEnv(conf)

	return conf, err
}
