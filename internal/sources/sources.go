package sources

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

type OptionsT struct {
	Name string
	Type string
	Path string

	Raw  string
	File string
	Git  OptionsGitT
	K8s  OptionsK8sT
}

type OptionsGitT struct {
	SshKeyFilepath string
	Url            string
	Branch         string
	Filepath       string
}

type OptionsK8sT struct {
	InCluster      bool
	ConfigFilepath string
	MasterUrl      string
	Kind           string
	Namespace      string
	Name           string
	Key            string
}

func GetSource(ops OptionsT) (SourceT, error) {
	switch ops.Type {
	case TypeRAW:
		{
			return NewRawSource(ops)
		}
	case TypeFILE:
		{
			return NewFileSource(ops)
		}
	case TypeGIT:
		{
			return NewGitSource(ops)
		}
	case TypeK8S:
		{
			return NewK8sSource(ops)
		}
	}
	return nil, nil
}
