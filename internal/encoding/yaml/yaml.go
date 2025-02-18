package yaml

import (
	"gopkg.in/yaml.v3"
)

// ----------------------------------------------------------------
// Decode/Encode/Merge YAML data structure
// ----------------------------------------------------------------

type YamlT struct {
}

func (e *YamlT) DecodeConfig(cfgBytes []byte) (cfg map[string]any, err error) {
	err = yaml.Unmarshal(cfgBytes, &cfg)
	return cfg, err
}

func (e *YamlT) EncodeConfig(cfg map[string]any) (configBytes []byte, err error) {
	return yaml.Marshal(cfg)
}

func (e *YamlT) MergeConfigs(dst, src map[string]any) error {
	mergeYamlObjects(dst, src)
	return nil
}
