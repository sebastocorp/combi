package raw

import (
	"os"
	"path/filepath"

	"combi/api/v1alpha3"
)

type RawSourceT struct {
	filepath string
}

func NewRawSource(srcConf v1alpha3.SourceConfigT, srcpath string) (s *RawSourceT, err error) {
	s = &RawSourceT{}

	s.filepath = filepath.Join(srcpath, "rawconfig.txt")

	err = os.WriteFile(s.filepath, []byte(srcConf.Raw), 0777)

	return s, err
}

func (s *RawSourceT) SyncConfig() (bool, error) {
	return false, nil
}

func (s *RawSourceT) GetConfig() (config []byte, err error) {
	config, err = os.ReadFile(s.filepath)

	return config, err
}
