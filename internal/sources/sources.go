package sources

import (
	"combi/api/v1alpha3"
	"combi/internal/config"
	"combi/internal/sources/raw"
)

type SourceT interface {
	SyncConfig() (bool, error)
	GetConfig() ([]byte, error)
}

func GetSource(srcConf v1alpha3.SourceConfigT, srcpath string) (SourceT, error) {
	switch srcConf.Type {
	case config.ConfigSourceTypeValueRAW:
		{
			return raw.NewRawSource(srcConf, srcpath)
		}
	case config.ConfigSourceTypeValueFILE:
		{
		}
	case config.ConfigSourceTypeValueGIT:
		{
		}
	case config.ConfigSourceTypeValueK8S:
		{
		}
	}
	return nil, nil
}
