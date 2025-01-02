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

func (e *JsonT) DecodeConfigBytes(configBytes []byte) (cfg map[string]any, err error) {
	if ok, err := regexp.Match("^[ ]*$", configBytes); ok {
		if err != nil {
			return cfg, err
		}
		configBytes = []byte("{}")
	}
	err = json.Unmarshal(configBytes, &cfg)
	return cfg, err
}

// Encode functions

func (e *JsonT) EncodeConfigString(cfg map[string]any) (configStr string) {
	configBytes, _ := json.MarshalIndent(cfg, "", "  ")
	configStr = string(configBytes)
	return configStr
}
