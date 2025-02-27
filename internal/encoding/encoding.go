package encoding

import (
	"combi/internal/encoding/json"
	"combi/internal/encoding/libconfig"
	"combi/internal/encoding/nginx"
	"combi/internal/encoding/yaml"
)

const (
	KindJSON      = "JSON"
	KindYAML      = "YAML"
	KindNGINX     = "NGINX"
	KindLIBCONFIG = "LIBCONFIG"
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
		KindJSON:      &json.JsonT{},
		KindYAML:      &yaml.YamlT{},
		KindLIBCONFIG: &libconfig.LibconfigT{},
		KindNGINX:     &nginx.NginxT{},
	}
	return encoders[encType]
}
