package raw

import (
	"os"
	"path/filepath"

	"combi/api/v1alpha3"
	"combi/internal/config"
)

type RawSourceT struct {
	storConfig string
}

func NewRawSource(srcConf v1alpha3.SourceConfigT, srcpath string) (s *RawSourceT, err error) {
	s = &RawSourceT{
		storConfig: filepath.Join(srcpath, "config.raw.txt"),
	}

	content := config.ExpandEnv([]byte(srcConf.Raw))
	err = os.WriteFile(s.storConfig, content, 0777)

	return s, err
}

func (s *RawSourceT) SyncConfig() (bool, error) {
	return false, nil
}

func (s *RawSourceT) GetConfig() (config []byte, err error) {
	config, err = os.ReadFile(s.storConfig)

	return config, err
}
