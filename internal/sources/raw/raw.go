package raw

import (
	"os"
	"path/filepath"

	"combi/api/v1alpha4"
	"combi/internal/config"
)

type RawSourceT struct {
	name       string
	storConfig string
}

func NewRawSource(srcConf v1alpha4.SourceConfigT, srcpath string) (s *RawSourceT, err error) {
	s = &RawSourceT{
		name:       srcConf.Name,
		storConfig: filepath.Join(srcpath, "config.raw.txt"),
	}

	content := config.ExpandEnv([]byte(srcConf.Raw))
	err = os.WriteFile(s.storConfig, content, 0777)

	return s, err
}
func (s *RawSourceT) GetName() string {
	return s.name
}

func (s *RawSourceT) SyncConfig() (bool, error) {
	return false, nil
}

func (s *RawSourceT) GetConfig() (config []byte, err error) {
	config, err = os.ReadFile(s.storConfig)

	return config, err
}
