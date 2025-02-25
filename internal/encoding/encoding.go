package encoding

import (
	"combi/internal/config"
	"combi/internal/encoding/json"
	"combi/internal/encoding/libconfig"
	"combi/internal/encoding/nginx"
	"combi/internal/encoding/yaml"
)

type EncoderT interface {
	// Encode/Decode configurations
	DecodeConfig([]byte) (map[string]any, error)
	EncodeConfig(map[string]any) ([]byte, error)

	// Merge configurations
	MergeConfigs(dst map[string]any, src map[string]any) error
}

func GetEncoder(encType string) EncoderT {
	encoders := map[string]EncoderT{
		config.ConfigKindValueJSON:      &json.JsonT{},
		config.ConfigKindValueYAML:      &yaml.YamlT{},
		config.ConfigKindValueLIBCONFIG: &libconfig.LibconfigT{},
		config.ConfigKindValueNGINX:     &nginx.NginxT{},
	}
	return encoders[encType]
}
