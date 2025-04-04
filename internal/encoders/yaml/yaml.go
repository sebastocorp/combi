package yaml

import (
	"gopkg.in/yaml.v3"
)

// ----------------------------------------------------------------
// Decode/Encode/Merge YAML data structure
// ----------------------------------------------------------------

type YamlT struct {
}

func (e *YamlT) Decode(cfgBytes []byte) (cfg map[string]any, err error) {
	err = yaml.Unmarshal(cfgBytes, &cfg)
	return cfg, err
}

func (e *YamlT) Encode(cfg map[string]any) (configBytes []byte, err error) {
	return yaml.Marshal(cfg)
}

func (e *YamlT) Merge(dst, src map[string]any) error {
	mergeYamlObjects(dst, src)
	return nil
}
