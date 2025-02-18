package libconfig


func mergeSettings(dst map[string]any, src map[string]any) {
	for srck := range src {
		switch src[srck].(type) {
		case string, []any:
			{
				dst[srck] = src[srck]
			}
		case map[string]any:
			{
				if _, ok := dst[srck]; ok {
					mergeSettings(dst[srck].(map[string]any), src[srck].(map[string]any))
				} else {
					dst[srck] = src[srck]
				}
			}
		}
	}
}
