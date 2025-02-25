package libconfig

func mergeSettings(dst map[string]any, src map[string]any) {
	for srck := range src {
		if _, ok := dst[srck]; !ok {
			dst[srck] = src[srck]
			continue
		}

		switch src[srck].(type) {
		case string, []any:
			{
				dst[srck] = src[srck]
			}
		case map[string]any:
			{
				mergeSettings(dst[srck].(map[string]any), src[srck].(map[string]any))
			}
		}
	}
}
