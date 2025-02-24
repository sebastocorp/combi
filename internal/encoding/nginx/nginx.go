package nginx

// ----------------------------------------------------------------
// Decode/Encode/Merge NGINX data structure
// ----------------------------------------------------------------

type NginxT struct {
}

func (e *NginxT) DecodeConfig(cfgBytes []byte) (cfg map[string]any, err error) {
	ts, err := tokenize(cfgBytes)
	if err != nil {
		return nil, err
	}

	cfg = make(map[string]any)
	err = decodeConfigTokens(ts, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (e *NginxT) EncodeConfig(cfg map[string]any) (result []byte, err error) {
	result = []byte(encodeConfig(cfg))
	return result, err
}

func (e *NginxT) MergeConfigs(dst map[string]any, src map[string]any) error {
	mergeConfig(dst, src)
	return nil
}
