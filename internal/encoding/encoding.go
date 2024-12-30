package encoding

import (
	"combi/internal/config"
	"combi/internal/encoding/json"
	"combi/internal/encoding/libconfig"
	"combi/internal/encoding/nginx"
)

type EncoderT interface {
	// Encode/Decode configurations
	DecodeConfig(filepath string) (err error)
	DecodeConfigBytes(configBytes []byte) (err error)
	EncodeConfigString() (configStr string)

	// Merge configurations
	MergeConfigs(source interface{})
	GetConfigStruct() (config interface{})

	// Transform configurations
	ConfigToMap() (configMap map[string]interface{})
}

func GetEncoders() (encoders map[string]EncoderT) {
	encoders = map[string]EncoderT{
		config.ConfigKindValueJSON:      &json.JsonT{},
		config.ConfigKindValueNGINX:     &nginx.NginxT{},
		config.ConfigKindValueLIBCONFIG: &libconfig.LibconfigT{},
	}
	return encoders
}

func GetEncoder(encType string) (encoder EncoderT) {
	encoders := map[string]EncoderT{
		config.ConfigKindValueJSON:      &json.JsonT{},
		config.ConfigKindValueNGINX:     &nginx.NginxT{},
		config.ConfigKindValueLIBCONFIG: &libconfig.LibconfigT{},
	}
	return encoders[encType]
}
