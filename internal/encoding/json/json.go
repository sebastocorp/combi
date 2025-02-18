package json

import (
	"encoding/json"
	"regexp"
)

// ----------------------------------------------------------------
// Decode/Encode/Merge JSON data structure
// ----------------------------------------------------------------

type JsonT struct {
}

func (e *JsonT) DecodeConfig(configBytes []byte) (cfg map[string]any, err error) {
	if ok, err := regexp.Match(`^\s*$`, configBytes); ok {
		if err != nil {
			return cfg, err
		}
		configBytes = []byte("{}")
	}
	err = json.Unmarshal(configBytes, &cfg)
	return cfg, err
}

func (e *JsonT) EncodeConfig(cfg map[string]any) (configBytes []byte, err error) {
	return json.MarshalIndent(cfg, "", "  ")
}

func (e *JsonT) MergeConfigs(dst, src map[string]any) error {
	mergeJsonObjects(dst, src)
	return nil
}
