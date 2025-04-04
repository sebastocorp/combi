package nginx

import "fmt"

func encodeConfig(cfg map[string]any) (result string) {
	result = encodeConfigIndent(cfg, 0)
	return result
}

func encodeConfigIndent(cfg map[string]any, indent int) (result string) {
	indentStr := ""
	for i := 0; i < indent; i++ {
		indentStr += "    "
	}

	resultDirectives := ""
	resultBocks := ""
	for cfgk := range cfg {
		if cfgk == nginxBlockPatternKey {
			fmt.Printf("skip block pattern key\n")
			continue
		}

		switch cfg[cfgk].(type) {
		case string:
			{
				resultDirectives += indentStr + cfgk + " " + cfg[cfgk].(string) + ";\n"
			}

		case []any:
			{
				switch cfg[cfgk].([]any)[0].(type) {
				case string:
					{
						for i := range cfg[cfgk].([]any) {
							resultDirectives += indentStr + cfgk + " " + cfg[cfgk].([]any)[i].(string) + ";\n"
						}
					}
				case map[string]any:
					{
						for i := range cfg[cfgk].([]any) {
							blockPattern := ""
							if cfg[cfgk].([]any)[i].(map[string]any)[nginxBlockPatternKey] != nil {
								blockPattern = " " + cfg[cfgk].([]any)[i].(map[string]any)[nginxBlockPatternKey].(string)
							}
							resultBocks += indentStr + cfgk + blockPattern + " {\n" + encodeConfigIndent(cfg[cfgk].([]any)[i].(map[string]any), indent+1) + indentStr + "}\n"
						}
					}
				}
			}
		case map[string]any:
			{
				blockPattern := ""
				if cfg[cfgk].(map[string]any)[nginxBlockPatternKey] != nil {
					blockPattern = " " + cfg[cfgk].(map[string]any)[nginxBlockPatternKey].(string)
				}
				resultBocks += indentStr + cfgk + blockPattern + " {\n" + encodeConfigIndent(cfg[cfgk].(map[string]any), indent+1) + indentStr + "}\n"
			}
		}
	}
	result = resultDirectives + resultBocks

	return result
}
