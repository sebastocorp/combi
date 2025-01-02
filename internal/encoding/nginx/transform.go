package nginx

import (
	"strings"
)

func (e *NginxT) ConfigToMap(cfgStruct *BlockContentT) (configMap map[string]any) {
	configMap = make(map[string]any)
	blockContentToMap(cfgStruct, configMap)
	return configMap
}

// func (e *NginxT) ConfigToMap() (configMap map[string]any) {
// 	configMap = make(map[string]any)
// 	blockContentToMap(&e.ConfigStruct, configMap)
// 	return configMap
// }

func blockContentToMap(blockContent *BlockContentT, configMap map[string]any) {
	for _, val := range blockContent.Directives {
		keyStr := val.Name
		if val.Param != "" {
			keyStr = strings.Join([]string{val.Name, "[", val.Param, "]"}, "")
		}
		configMap[keyStr] = val.Value
	}

	for _, val := range blockContent.Blocks {
		keyStr := val.Name
		if val.Params != "" {
			keyStr = strings.Join([]string{val.Name, "[", val.Params, "]"}, "")
		}
		configMap[keyStr] = make(map[string]any)
		blockContentToMap(&val.BlockContent, configMap[keyStr].(map[string]any))
	}
}
