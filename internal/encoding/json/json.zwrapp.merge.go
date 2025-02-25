package json

func mergeJsonObjects(dst, src map[string]any) {
	for srck := range src {
		if _, ok := dst[srck]; !ok {
			dst[srck] = src[srck]
			continue
		}

		switch dst[srck].(type) {
		case float64, string, bool, []any, nil:
			dst[srck] = src[srck]
		case map[string]any:
			{
				mergeJsonObjects(dst[srck].(map[string]any), src[srck].(map[string]any))
			}
		}
	}
}
