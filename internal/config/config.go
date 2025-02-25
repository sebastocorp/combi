package config

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"time"

	"combi/api/v1alpha4"

	"gopkg.in/yaml.v3"
)

const (
	ConfigKindValueJSON      = "JSON"
	ConfigKindValueYAML      = "YAML"
	ConfigKindValueNGINX     = "NGINX"
	ConfigKindValueLIBCONFIG = "LIBCONFIG"

	ConfigOnValueSUCCESS = "SUCCESS"
	ConfigOnValueFAILURE = "FAILURE"

	ConfigSourceTypeValueRAW  = "RAW"
	ConfigSourceTypeValueFILE = "FILE"
	ConfigSourceTypeValueGIT  = "GIT"
	ConfigSourceTypeValueK8S  = "K8S"
)

// ParseConfig TODO
func ParseConfig(cfgBytes []byte) (cfg v1alpha4.CombiConfigT, err error) {
	cfgBytes = ExpandEnv(cfgBytes)

	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return cfg, err
	}

	err = checkConfig(&cfg)

	return cfg, err
}

// ExpandEnv TODO
func ExpandEnv(input []byte) []byte {
	re := regexp.MustCompile(`\${ENV:([A-Za-z_][A-Za-z0-9_]*)}\$`)
	result := re.ReplaceAllFunc(input, func(match []byte) []byte {
		key := re.FindSubmatch(match)[1]
		if value, exists := os.LookupEnv(string(key)); exists {
			return []byte(value)
		}
		return match
	})

	return result
}

// checkConfig TODO
func checkConfig(cfg *v1alpha4.CombiConfigT) error {
	spacesRegex := regexp.MustCompile(`\s`)

	configKindValues := []string{
		ConfigKindValueJSON,
		ConfigKindValueNGINX,
		ConfigKindValueLIBCONFIG,
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

	if cfg.Settings.TmpObjs.Path == "" {
		cfg.Settings.TmpObjs.Path = "/tmp/combi"
	}

	if cfg.Settings.TmpObjs.Mode == 0 {
		cfg.Settings.TmpObjs.Mode = 0777
	}

	//------------------------------
	// Sources
	//------------------------------

	namesCount := map[string]int{}
	srcTypeValues := []string{
		ConfigSourceTypeValueRAW,
		ConfigSourceTypeValueFILE,
		ConfigSourceTypeValueGIT,
		ConfigSourceTypeValueK8S,
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
		ConfigOnValueSUCCESS,
		ConfigOnValueFAILURE,
	}
	for ai := range cfg.Behavior.Actions {
		if matched := spacesRegex.MatchString(cfg.Behavior.Actions[ai].Name); matched {
			return fmt.Errorf("action name '%s' contains spaces", cfg.Behavior.Actions[ai].Name)
		}

		if !slices.Contains(onValues, cfg.Behavior.Actions[ai].On) {
			return fmt.Errorf("action '%s' on field must be one of this %v", cfg.Behavior.Actions[ai].Name, onValues)
		}

		if _, ok := namesCount[cfg.Behavior.Actions[ai].Name]; ok {
			return fmt.Errorf("action '%s' is duplicated", cfg.Behavior.Actions[ai].Name)
		}
		namesCount[cfg.Behavior.Actions[ai].Name] += 1
	}

	return nil
}
