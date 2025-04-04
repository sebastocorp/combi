package nginx

func mergeConfig(dst map[string]any, src map[string]any) {
	for srck := range src {
		if _, ok := dst[srck]; !ok {
			dst[srck] = src[srck]
			continue
		}

		switch src[srck].(type) {
		case string:
			{
				dst[srck] = src[srck]
			}
		case []any:
			{
				dst[srck] = src[srck]
			}
		case map[string]any:
			{
				mergeConfig(dst[srck].(map[string]any), src[srck].(map[string]any))
			}
		}
	}
}
