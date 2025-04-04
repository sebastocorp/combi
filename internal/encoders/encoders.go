package encoders

import (
	"combi/internal/encoders/json"
	"combi/internal/encoders/libconfig"
	"combi/internal/encoders/nginx"
	"combi/internal/encoders/yaml"
)

const (
	TypeJSON      = "JSON"
	TypeYAML      = "YAML"
	TypeNGINX     = "NGINX"
	TypeLIBCONFIG = "LIBCONFIG"
)

var (
	// EncoderTypes is a list of supported encoder types
	Encoders = map[string]EncoderT{
		TypeJSON:      &json.JsonT{},
		TypeYAML:      &yaml.YamlT{},
		TypeLIBCONFIG: &libconfig.LibconfigT{},
		TypeNGINX:     &nginx.NginxT{},
	}
)

type EncoderT interface {
	// Encode/Decode configurations
	Decode([]byte) (map[string]any, error)
	Encode(map[string]any) ([]byte, error)

	// Merge configurations
	Merge(dst map[string]any, src map[string]any) error
}
