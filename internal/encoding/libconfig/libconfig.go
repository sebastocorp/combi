package libconfig

// ----------------------------------------------------------------
// Decode/Encode/Merge Libconfig data structure
// ----------------------------------------------------------------

type LibconfigT struct {
}

func (e *LibconfigT) DecodeConfig(cfgBytes []byte) (cfg map[string]any, err error) {
	ts, err := tokenize(cfgBytes)
	if err != nil {
		return nil, err
	}

	cfg = make(map[string]any)
	err = decodeSettingsTokens(ts, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (e *LibconfigT) EncodeConfig(cfg map[string]any) (result []byte, err error) {
	result = []byte(encodeSettings(cfg))
	return result, err
}

func (e *LibconfigT) MergeConfigs(dst map[string]any, src map[string]any) error {
	mergeSettings(dst, src)
	return nil
}
