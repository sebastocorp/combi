package libconfig

import (
	"os"
	"regexp"
)

type LibconfigT struct {
	ConfigStruct map[string]any
}

// ----------------------------------------------------------------
// Decode/Encode LIBCONFIG data structure
// ----------------------------------------------------------------

// Decode functions

func (e *LibconfigT) DecodeConfig(filepath string) (err error) {
	configBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	configBytes = regexp.MustCompile(`#[^\n]*`).ReplaceAll(configBytes, []byte(""))

	err = e.DecodeConfigBytes(configBytes)
	return err
}

func (e *LibconfigT) DecodeConfigBytes(configBytes []byte) (err error) {
	configStr := string(configBytes)
	configStr = regexp.MustCompile(`#[^\n]*`).ReplaceAllString(configStr, "")

	err = e.parseLibconfigString(configStr)
	return err
}

// Encode functions

func (e *LibconfigT) EncodeConfigString() (configStr string) {
	configStr += encodeConfigSettingString(e.ConfigStruct, 0)
	return configStr
}

func encodeConfigSettingString(settings map[string]any, indent int) (configStr string) {
	var indentStr string
	for i := 0; i < indent; i++ {
		indentStr += "  "
	}

	// Encode settings with primitive
	for name, value := range settings {
		switch value.(type) {
		case string:
			configStr += indentStr + name + "=" + value.(string) + ",\n"
		case []string:
			configStr += indentStr + name + "=" + "\n"
			configStr += encodeConfigArrayString(value.([]string), indent)
			configStr += ",\n"
		case []any:
			configStr += indentStr + name + "=" + "\n"
			configStr += encodeConfigListString(value.([]any), indent)
			configStr += ",\n"
		case map[string]any:
			configStr += indentStr + name + "=" + "\n"
			configStr += encodeConfigGroupString(value.(map[string]any), indent)
			configStr += ",\n"
		default:
			// logger.Log.Debugf("invalid libconfig type\n")
		}
	}

	return configStr
}

func encodeConfigArrayString(array []string, indent int) (configStr string) {
	var indentStr string
	for i := 0; i < indent; i++ {
		indentStr += "  "
	}

	configStr += indentStr + "[\n" + indentStr + "  "
	for _, primitive := range array {
		configStr += primitive + ", "
	}
	configStr += "\n" + indentStr + "]"
	return configStr
}

func encodeConfigGroupString(group map[string]any, indent int) (configStr string) {
	var indentStr string
	for i := 0; i < indent; i++ {
		indentStr += "  "
	}

	configStr += indentStr + "{\n"
	configStr += encodeConfigSettingString(group, indent+1)
	configStr += indentStr + "}"
	return configStr
}

func encodeConfigListString(list []any, indent int) (configStr string) {
	var indentStr string
	for i := 0; i < indent; i++ {
		indentStr += "  "
	}

	configStr += indentStr + "(\n"
	for index, value := range list {
		switch value.(type) {
		case string:
			configStr += value.(string)
		case []string:
			configStr += encodeConfigArrayString(value.([]string), indent+1)
		case []any:
			configStr += encodeConfigListString(value.([]any), indent+1)
		case map[string]any:
			configStr += encodeConfigGroupString(value.(map[string]any), indent+1)
		default:
			// logger.Log.Debugf("invalid libconfig type\n")
		}

		if index < len(list)-1 {
			configStr += ",\n"
		}
	}

	configStr += "\n" + indentStr + ")"
	return configStr
}
