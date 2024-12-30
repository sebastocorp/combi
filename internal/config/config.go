package config

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"time"

	"combi/api/v1alpha3"

	"gopkg.in/yaml.v3"
)

const (
	ConfigKindValueJSON      = "JSON"
	ConfigKindValueNGINX     = "NGINX"
	ConfigKindValueLIBCONFIG = "LIBCONFIG"

	ConfigConditionResultValueSUCCESS = "SUCCESS"
	ConfigConditionResultValueFAILURE = "FAILURE"

	ConfigSourceTypeValueRAW  = "RAW"
	ConfigSourceTypeValueFILE = "FILE"
	ConfigSourceTypeValueGIT  = "GIT"
	ConfigSourceTypeValueK8S  = "K8S"
)

// ParseConfig TODO
func ParseConfig(fileBytes []byte) (config v1alpha3.CombiConfigT, err error) {
	fileBytes = ExpandEnv(fileBytes)

	err = yaml.Unmarshal(fileBytes, &config)
	if err != nil {
		return config, err
	}

	err = checkConfig(&config)

	return config, err
}

// expandEnv TODO
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
func checkConfig(config *v1alpha3.CombiConfigT) error {
	configKindValues := []string{
		ConfigKindValueJSON,
		ConfigKindValueNGINX,
		ConfigKindValueLIBCONFIG,
	}
	if !slices.Contains(configKindValues, config.Kind) {
		return fmt.Errorf("kind must be one of this %v", configKindValues)
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
	for _, sv := range config.Sources {
		if !slices.Contains(srcTypeValues, sv.Type) {
			return fmt.Errorf("source type '%s' must be one of this %v", sv.Name, srcTypeValues)
		}

		if _, ok := namesCount[sv.Name]; ok {
			return fmt.Errorf("source '%s' is duplicated", sv.Name)
		}
		namesCount[sv.Name] += 1
	}

	//------------------------------
	// Behavior
	//------------------------------

	if config.Behavior.SyncTime < 2*time.Second {
		return fmt.Errorf("behavior.syncTime must be at least 2 seconds")
	}

	if config.Behavior.Target.Path == "" {
		return fmt.Errorf("behavior.target.path must be set")
	}

	if config.Behavior.Target.File == "" {
		return fmt.Errorf("behavior.target.file must be set")
	}

	if config.Behavior.Target.Mode == 0 {
		config.Behavior.Target.Mode = 0777
	}

	if config.Behavior.TmpObjs.Path == "" {
		config.Behavior.TmpObjs.Path = "/tmp/combi"
	}

	if config.Behavior.TmpObjs.Mode == 0 {
		config.Behavior.TmpObjs.Mode = 0777
	}

	// conditions check

	namesCount = map[string]int{}
	for _, cv := range config.Behavior.Conditions {
		if _, ok := namesCount[cv.Name]; ok {
			return fmt.Errorf("source '%s' is duplicated", cv.Name)
		}
		namesCount[cv.Name] += 1
	}

	// actions check

	namesCount = map[string]int{}
	conditionResultValues := []string{
		ConfigConditionResultValueSUCCESS,
		ConfigConditionResultValueFAILURE,
	}
	for _, av := range config.Behavior.Actions {
		if !slices.Contains(conditionResultValues, av.ConditionResult) {
			return fmt.Errorf("action '%s' conditionResult must be one of this %v", av.Name, conditionResultValues)
		}

		if _, ok := namesCount[av.Name]; ok {
			return fmt.Errorf("action '%s' is duplicated", av.Name)
		}
		namesCount[av.Name] += 1
	}

	return nil
}
