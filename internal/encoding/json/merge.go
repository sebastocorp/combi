package json

// ----------------------------------------------------------------
// Merge JSON data structure
// ----------------------------------------------------------------

func (e *JsonT) MergeConfigs(dst map[string]any, src map[string]any) error {
	mergeJsonObjects(dst, src)
	return nil
}

func mergeJsonObjects(destination, source map[string]any) {
	for srcKey, srcVal := range source {

		if _, ok := destination[srcKey]; !ok {
			destination[srcKey] = srcVal
			continue
		}

		switch destination[srcKey].(type) {
		case float64, string, bool, nil:
			destination[srcKey] = srcVal
		case []any:
			mergeJsonArray(destination[srcKey].([]any), srcVal.([]any))
		case map[string]any:
			mergeJsonObjects(destination[srcKey].(map[string]any), srcVal.(map[string]any))
		default:
			// logger.Log.Debugf("invalid json type\n")
		}
	}
}

func mergeJsonArray(destination, source []interface{}) {
	gap := len(source) - len(destination)
	if gap > 0 {
		for i := 0; i < gap; i++ {
			destination = append(destination, nil)
		}
	}
	for srcIndex, srcVal := range source {
		switch srcVal.(type) {
		case float64, string, bool, nil:
			destination[srcIndex] = srcVal
		case []any:
			{
				if destination[srcIndex] == nil {
					destination[srcIndex] = []any{}
				}
				mergeJsonArray(destination[srcIndex].([]any), srcVal.([]any))
			}
		case map[string]any:
			{
				if destination[srcIndex] == nil {
					destination[srcIndex] = map[string]any{}
				}
				mergeJsonObjects(destination[srcIndex].(map[string]any), srcVal.(map[string]any))
			}
		default:
			// logger.Log.Debugf("invalid json type\n")
		}
	}
}
