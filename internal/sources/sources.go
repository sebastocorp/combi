package sources

import (
	"combi/api/v1alpha4"
	"combi/internal/config"
	"combi/internal/sources/file"
	"combi/internal/sources/git"
	"combi/internal/sources/k8s"
	"combi/internal/sources/raw"
)

type SourceT interface {
	GetName() string
	SyncConfig() (bool, error)
	GetConfig() ([]byte, error)
}

func GetSource(srcCfg v1alpha4.SourceConfigT, srcpath string) (SourceT, error) {
	switch srcCfg.Type {
	case config.ConfigSourceTypeValueRAW:
		{
			return raw.NewRawSource(srcCfg, srcpath)
		}
	case config.ConfigSourceTypeValueFILE:
		{
			return file.NewFileSource(srcCfg, srcpath)
		}
	case config.ConfigSourceTypeValueGIT:
		{
			return git.NewGitSource(srcCfg, srcpath)
		}
	case config.ConfigSourceTypeValueK8S:
		{
			return k8s.NewK8sSource(srcCfg, srcpath)
		}
	}
	return nil, nil
}
