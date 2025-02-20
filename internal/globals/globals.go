package globals

import "strings"

const (
	LogKeyService         = "service"
	LogKeySourceName      = "sourceName"
	LogKeyCondition       = "condition"
	LogKeyConditionResult = "conditionResult"
	LogKeyAction          = "action"
	LogKeyActionStdout    = "actionStdout"
	LogKeyActionStderr    = "actionStderr"
	LogKeyError           = "error"

	LogValueService                = "combi"
	LogValueDefaultStr             = "none"
	LogValueConditionResultFAIL    = "FAIL"
	LogValueConditionResultSUCCESS = "SUCCESS"
)

func GetLogCommonFields() map[string]any {
	return map[string]any{
		LogKeyService: LogValueService,
	}
}

// CopyMap return a map that is a real copy of the original
// Ref: https://go.dev/blog/maps
func CopyMap(src map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(src))
	for k, v := range src {
		m[k] = v
	}
	return m
}

// SplitCommaSeparatedValues get a list of strings and return a new list
// where each element containing commas is divided in separated elements
func SplitCommaSeparatedValues(input []string) []string {
	var result []string
	for _, item := range input {
		parts := strings.Split(item, ",")
		result = append(result, parts...)
	}
	return result
}
