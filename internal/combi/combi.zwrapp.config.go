package combi

import (
	"fmt"
	"regexp"
	"slices"
	"time"

	"combi/api/v1alpha4"
	"combi/internal/combi/actionset"
	"combi/internal/combi/conditionset"
	"combi/internal/encoding"
	"combi/internal/sources"
	"combi/internal/utils"

	"gopkg.in/yaml.v3"
)

type ApiVersionConfigT struct {
	ApiVersion string `yaml:"apiVersion"`
}

// parseConfig TODO
func parseConfig(cfgBytes []byte) (cfg any, err error) {
	cfgBytes = utils.ExpandEnv(cfgBytes)

	avc := ApiVersionConfigT{}
	err = yaml.Unmarshal(cfgBytes, &avc)
	if err != nil {
		return cfg, err
	}

	switch avc.ApiVersion {
	case "combi/v1alpha4":
		{
			cfg, err = v1alpha4Parse(cfgBytes)
		}
	default:
		{
			return cfg, fmt.Errorf("unsupported apiVersion '%s'", avc.ApiVersion)
		}
	}

	return cfg, err
}

func v1alpha4Parse(cfgBytes []byte) (cfg v1alpha4.CombiConfigT, err error) {
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return cfg, err
	}

	err = v1alpha4Check(&cfg)
	return cfg, err
}

// checkConfig TODO
func v1alpha4Check(cfg *v1alpha4.CombiConfigT) error {
	spacesRegex := regexp.MustCompile(`\s`)

	configKindValues := []string{
		encoding.KindJSON,
		encoding.KindYAML,
		encoding.KindNGINX,
		encoding.KindLIBCONFIG,
	}
	if !slices.Contains(configKindValues, cfg.Kind) {
		return fmt.Errorf("kind field must be one of this %v", configKindValues)
	}

	//------------------------------
	// Settings
	//------------------------------

	if cfg.Settings.SyncTime < 2*time.Second {
		return fmt.Errorf("settings.syncTime field must be at least 2 seconds")
	}

	if cfg.Settings.Target.Path == "" {
		return fmt.Errorf("settings.target.path field must be set")
	}

	if cfg.Settings.Target.File == "" {
		return fmt.Errorf("settings.target.file field must be set")
	}

	if cfg.Settings.Target.Mode == 0 {
		cfg.Settings.Target.Mode = 0777
	}

	if cfg.Settings.TmpFiles.Path == "" {
		cfg.Settings.TmpFiles.Path = "/tmp/combi"
	}

	if cfg.Settings.TmpFiles.Mode == 0 {
		cfg.Settings.TmpFiles.Mode = 0777
	}

	//------------------------------
	// Sources
	//------------------------------

	namesCount := map[string]int{}
	srcTypeValues := []string{
		sources.TypeFILE,
		sources.TypeRAW,
		sources.TypeGIT,
		sources.TypeK8S,
	}
	for si := range cfg.Sources {
		if !slices.Contains(srcTypeValues, cfg.Sources[si].Type) {
			return fmt.Errorf("source '%s' type field must be one of this %v", cfg.Sources[si].Name, srcTypeValues)
		}

		if matched := spacesRegex.MatchString(cfg.Sources[si].Name); matched {
			return fmt.Errorf("source name '%s' contains spaces", cfg.Sources[si].Name)
		}

		if _, ok := namesCount[cfg.Sources[si].Name]; ok {
			return fmt.Errorf("source '%s' is duplicated", cfg.Sources[si].Name)
		}
		namesCount[cfg.Sources[si].Name] += 1
	}

	//------------------------------
	// Behavior
	//------------------------------

	// conditions check

	namesCount = map[string]int{}
	for ci := range cfg.Behavior.Conditions {
		if matched := spacesRegex.MatchString(cfg.Behavior.Conditions[ci].Name); matched {
			return fmt.Errorf("condition name '%s' contains spaces", cfg.Behavior.Conditions[ci].Name)
		}

		if _, ok := namesCount[cfg.Behavior.Conditions[ci].Name]; ok {
			return fmt.Errorf("source '%s' is duplicated", cfg.Behavior.Conditions[ci].Name)
		}
		namesCount[cfg.Behavior.Conditions[ci].Name] += 1
	}

	// actions check

	namesCount = map[string]int{}
	onValues := []string{
		conditionset.StatusSuccess,
		conditionset.StatusFail,
	}
	inValues := []string{
		actionset.TypeLOCAL,
		actionset.TypeK8S,
	}
	for ai := range cfg.Behavior.Actions {
		if matched := spacesRegex.MatchString(cfg.Behavior.Actions[ai].Name); matched {
			return fmt.Errorf("action name '%s' contains spaces", cfg.Behavior.Actions[ai].Name)
		}

		if !slices.Contains(onValues, cfg.Behavior.Actions[ai].On) {
			return fmt.Errorf("action '%s' on field must be one of this %v", cfg.Behavior.Actions[ai].Name, onValues)
		}

		if cfg.Behavior.Actions[ai].In == "" {
			cfg.Behavior.Actions[ai].In = actionset.TypeLOCAL
		}
		if !slices.Contains(inValues, cfg.Behavior.Actions[ai].In) {
			return fmt.Errorf("action '%s' in field must be one of this %v", cfg.Behavior.Actions[ai].Name, inValues)
		}

		if _, ok := namesCount[cfg.Behavior.Actions[ai].Name]; ok {
			return fmt.Errorf("action '%s' is duplicated", cfg.Behavior.Actions[ai].Name)
		}
		namesCount[cfg.Behavior.Actions[ai].Name] += 1
	}

	return nil
}
