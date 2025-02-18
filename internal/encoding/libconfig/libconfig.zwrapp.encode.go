package libconfig

func encodeSettings(cfg map[string]any) (result string) {
	resultSetts := ""
	resultArrs := ""
	resultLists := ""
	resultGroups := ""
	for cfgk := range cfg {
		switch cfg[cfgk].(type) {
		case string:
			{
				resultSetts += cfgk + "=" + cfg[cfgk].(string) + "\n"
			}

		case []any:
			{
				if isList(cfg[cfgk].([]any)) {
					resultLists += cfgk + "=\n" + encodeList(cfg[cfgk].([]any), 0) + "\n"
				} else {
					resultArrs += cfgk + "=" + encodeArray(cfg[cfgk].([]any)) + "\n"
				}
			}
		case map[string]any:
			{
				resultGroups += cfgk + "=\n" + encodeGroup(cfg[cfgk].(map[string]any), 0) + "\n"
			}
		}
	}
	result = resultSetts + resultArrs + resultGroups + resultLists

	return result
}

func encodeArray(cfg []any) (result string) {
	result += "[ "
	for i := range cfg {
		result += cfg[i].(string)
		if i < len(cfg)-1 {
			result += ", "
		}
	}
	result += " ]"

	return result
}

func isList(cfg []any) bool {
	isList := false
	for i := range cfg {
		switch cfg[i].(type) {
		case string:
		default:
			isList = true
		}
	}

	return isList
}

func encodeList(cfg []any, indentn int) (result string) {
	indent := ""
	for i := 0; i < indentn; i++ {
		indent += "  "
	}
	indentIn := indent + "  "

	result += indent + "(\n"
	for i := range cfg {
		switch cfg[i].(type) {
		case string:
			{
				result += indentIn + cfg[i].(string)
			}
		case []any:
			{
				if isList(cfg[i].([]any)) {
					result += encodeList(cfg[i].([]any), indentn+1)
				} else {
					result += indentIn + encodeArray(cfg[i].([]any))
				}
			}
		case map[string]any:
			{
				result += encodeGroup(cfg[i].(map[string]any), indentn+1)
			}
		}
		if i < len(cfg)-1 {
			result += ",\n"
		}
	}
	result += "\n" + indent + ")"

	return result
}

func encodeGroup(cfg map[string]any, indentn int) (result string) {
	indent := ""
	for i := 0; i < indentn; i++ {
		indent += "  "
	}
	indentIn := indent + "  "

	result += indent + "{\n"
	index := 0
	for cfgk := range cfg {
		result += indentIn + cfgk + "="
		switch cfg[cfgk].(type) {
		case string:
			{
				result += cfg[cfgk].(string)
			}
		case []any:
			{
				if isList(cfg[cfgk].([]any)) {
					result += encodeList(cfg[cfgk].([]any), indentn+1)
				} else {
					result += encodeArray(cfg[cfgk].([]any))
				}
			}
		case map[string]any:
			{
				result += encodeGroup(cfg[cfgk].(map[string]any), indentn+1)
			}
		}
		if index < len(cfg)-1 {
			result += ",\n"
		}
		index++
	}
	result += "\n" + indent + "}"

	return result
}
