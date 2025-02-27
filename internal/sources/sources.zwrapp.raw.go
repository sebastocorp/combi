package sources

import (
	"combi/internal/utils"
	"os"
	"path/filepath"
)

type RawSourceT struct {
	name       string
	storConfig string
}

func NewRawSource(ops OptionsT) (s *RawSourceT, err error) {
	s = &RawSourceT{
		name:       ops.Name,
		storConfig: filepath.Join(ops.Path, "config.raw.txt"),
	}

	content := utils.ExpandEnv([]byte(ops.Raw))
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
