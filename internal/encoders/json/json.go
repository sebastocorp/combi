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

func (e *JsonT) Decode(configBytes []byte) (cfg map[string]any, err error) {
	if ok, err := regexp.Match(`^\s*$`, configBytes); ok {
		if err != nil {
			return cfg, err
		}
		configBytes = []byte("{}")
	}
	err = json.Unmarshal(configBytes, &cfg)
	return cfg, err
}

func (e *JsonT) Encode(cfg map[string]any) (configBytes []byte, err error) {
	return json.MarshalIndent(cfg, "", "  ")
}

func (e *JsonT) Merge(dst, src map[string]any) error {
	mergeJsonObjects(dst, src)
	return nil
}
