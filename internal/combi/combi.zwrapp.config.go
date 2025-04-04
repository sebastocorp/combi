package combi

import (
	"fmt"

	"combi/api/v1alpha5"
	"combi/internal/utils"

	"gopkg.in/yaml.v3"
)

type AVKT struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
}

// parseConfig TODO
func parseConfig(cfgBytes []byte) (cfg any, err error) {
	cfgBytes = utils.ExpandEnv(cfgBytes)

	avk := AVKT{}
	err = yaml.Unmarshal(cfgBytes, &avk)
	if err != nil {
		return cfg, err
	}

	switch avk.ApiVersion {
	case "combi/v1alpha5":
		{
			if avk.Kind != "Configuration" {
				return cfg, fmt.Errorf("not supported kind in apiVersion, must be 'Combi'")
			}
			cfg, err = v1alpha5Parse(cfgBytes)
		}
	default:
		{
			return cfg, fmt.Errorf("unsupported apiVersion '%s'", avk.ApiVersion)
		}
	}

	return cfg, err
}

// v1alpha5Parse TODO
func v1alpha5Parse(cfgBytes []byte) (cfg v1alpha5.CombiT, err error) {
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return cfg, err
	}

	err = v1alpha5Check(&cfg)
	return cfg, err
}

// v1alpha5Check TODO
func v1alpha5Check(cfg *v1alpha5.CombiT) error {
	return nil
}
