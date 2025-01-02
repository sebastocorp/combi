package encoding

import (
	"combi/internal/config"
	"combi/internal/encoding/json"
	"combi/internal/encoding/libconfig"
)

type EncoderT interface {
	// Encode/Decode configurations
	DecodeConfigBytes([]byte) (map[string]any, error)
	EncodeConfigString(map[string]any) string

	// Merge configurations
	MergeConfigs(dst map[string]any, src map[string]any)
}

// func GetEncoders() (encoders map[string]EncoderT) {
// 	encoders = map[string]EncoderT{
// 		config.ConfigKindValueJSON:      &json.JsonT{},
// 		config.ConfigKindValueNGINX:     &nginx.NginxT{},
// 		config.ConfigKindValueLIBCONFIG: &libconfig.LibconfigT{},
// 	}
// 	return encoders
// }

func GetEncoder(encType string) (encoder EncoderT) {
	encoders := map[string]EncoderT{
		config.ConfigKindValueJSON: &json.JsonT{},
		// config.ConfigKindValueNGINX:     &nginx.NginxT{},
		config.ConfigKindValueLIBCONFIG: &libconfig.LibconfigT{},
	}
	return encoders[encType]
}
