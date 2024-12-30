package file

import (
	"combi/api/v1alpha3"
	"combi/internal/config"
	"os"
	"path/filepath"
	"reflect"
)

type FileSourceT struct {
	srcFilepath     string
	currentFilepath string
}

func NewFileSource(srcConf v1alpha3.SourceConfigT, srcpath string) (s *FileSourceT, err error) {
	s = &FileSourceT{
		srcFilepath: srcConf.File,
	}

	var configBytes []byte
	configBytes, err = os.ReadFile(s.srcFilepath)
	if err != nil {
		return s, err
	}

	s.currentFilepath = filepath.Join(srcpath, filepath.Base(s.srcFilepath))
	err = os.WriteFile(s.currentFilepath, configBytes, 0777)
	if err != nil {
		return s, err
	}

	return s, err
}

func (s *FileSourceT) SyncConfig() (updated bool, err error) {
	syncBytes, err := os.ReadFile(s.srcFilepath)
	if err != nil {
		return updated, err
	}

	currentBytes, err := os.ReadFile(s.srcFilepath)
	if err != nil {
		return updated, err
	}

	if !reflect.DeepEqual(syncBytes, currentBytes) {
		updated = true
		err = os.WriteFile(s.currentFilepath, syncBytes, 0777)
		if err != nil {
			return updated, err
		}
	}

	return updated, err
}

func (s *FileSourceT) GetConfig() (conf []byte, err error) {
	if conf, err = os.ReadFile(s.currentFilepath); err != nil {
		return conf, err
	}

	conf = config.ExpandEnv(conf)

	return conf, err
}
