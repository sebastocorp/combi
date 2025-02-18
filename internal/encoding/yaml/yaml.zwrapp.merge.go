package yaml

func mergeYamlObjects(dst, src map[string]any) {
	for srck := range src {
		switch dst[srck].(type) {
		case float64, string, bool, []any, nil:
			dst[srck] = src[srck]
		case map[string]any:
			{
				if _, ok := dst[srck]; ok {
					mergeYamlObjects(dst[srck].(map[string]any), src[srck].(map[string]any))
				} else {
					dst[srck] = src[srck]
				}
			}
		}
	}
}
