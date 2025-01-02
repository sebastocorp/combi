package json

import (
	"encoding/json"
	"regexp"
)

type JsonT struct {
}

// ----------------------------------------------------------------
// Decode/Encode JSON data structure
// ----------------------------------------------------------------

// Decode functions

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

// Encode functions

func (e *JsonT) EncodeConfig(cfg map[string]any) (configBytes []byte, err error) {
	return json.MarshalIndent(cfg, "", "  ")
}
