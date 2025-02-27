package sources

import (
	"combi/api/v1alpha4"
	"combi/internal/sources/file"
	"combi/internal/sources/git"
	"combi/internal/sources/k8s"
	"combi/internal/sources/raw"
)

const (
	TypeRAW  = "RAW"
	TypeFILE = "FILE"
	TypeGIT  = "GIT"
	TypeK8S  = "K8S"
)

type SourceT interface {
	GetName() string
	SyncConfig() (bool, error)
	GetConfig() ([]byte, error)
}

func GetSource(srcCfg v1alpha4.SourceConfigT, srcpath string) (SourceT, error) {
	switch srcCfg.Type {
	case TypeRAW:
		{
			return raw.NewRawSource(srcCfg, srcpath)
		}
	case TypeFILE:
		{
			return file.NewFileSource(srcCfg, srcpath)
		}
	case TypeGIT:
		{
			return git.NewGitSource(srcCfg, srcpath)
		}
	case TypeK8S:
		{
			return k8s.NewK8sSource(srcCfg, srcpath)
		}
	}
	return nil, nil
}
