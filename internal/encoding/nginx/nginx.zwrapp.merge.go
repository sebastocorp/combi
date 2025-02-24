package nginx

func mergeConfig(dst map[string]any, src map[string]any) {
	for srck := range src {
		switch src[srck].(type) {
		case string, []any:
			{
				dst[srck] = src[srck]
			}
		case map[string]any:
			{
				if _, ok := dst[srck]; ok {
					mergeConfig(dst[srck].(map[string]any), src[srck].(map[string]any))
				} else {
					dst[srck] = src[srck]
				}
			}
		}
	}
}
