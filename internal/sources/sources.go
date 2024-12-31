package sources

import (
	"combi/api/v1alpha3"
	"combi/internal/config"
	"combi/internal/sources/file"
	"combi/internal/sources/git"
	"combi/internal/sources/k8s"
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
			return file.NewFileSource(srcConf, srcpath)
		}
	case config.ConfigSourceTypeValueGIT:
		{
			return git.NewGitSource(srcConf, srcpath)
		}
	case config.ConfigSourceTypeValueK8S:
		{
			return k8s.NewK8sSource(srcConf, srcpath)
		}
	}
	return nil, nil
}
