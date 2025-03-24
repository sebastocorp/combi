package v1alpha5

import "time"

type CombiT struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Conf       ConfT  `yaml:"conf"`
}

//--------------------------------------------------------------
// CONFIG
//--------------------------------------------------------------

type ConfT struct {
	SyncTime time.Duration `yaml:"syncTime"`
	Logger   LoggerT       `yaml:"logger"`
	TmpFiles TmpFilesT     `yaml:"tmpFiles"`

	Sources SourcesT `yaml:"sources"`
	Target  TargetT  `yaml:"target"`
}

type LoggerT struct {
	Level string `yaml:"level"`
}

type TmpFilesT struct {
	Path string `yaml:"path"`
	Mode uint32 `yaml:"mode"`
}

//--------------------------------------------------------------
// SOURCES CONFIG
//--------------------------------------------------------------

type SourcesT struct {
	Encoder     string        `yaml:"encoder"`
	Credentials []CredentialT `yaml:"credentials"`
	List        []SourceT     `yaml:"list"`
}

type CredentialT struct {
}

type SourceT struct {
	Name string     `yaml:"name"`
	Type string     `yaml:"type"` // values: RAW|FILE|GIT|K8S
	Raw  string     `yaml:"raw,omitempty"`
	File string     `yaml:"file,omitempty"`
	Git  SourceGitT `yaml:"git,omitempty"`
	K8s  SourceK8sT `yaml:"k8s,omitempty"`
}

type SourceGitT struct {
	SshUrl         string `yaml:"sshUrl"`
	SshKeyFilepath string `yaml:"sshKeyFilepath"`
	Branch         string `yaml:"branch"`
	Filepath       string `yaml:"filepath"`
}

type SourceK8sT struct {
	Context   SourceK8sContextConfigT `yaml:"context"`
	Kind      string                  `yaml:"kind"`
	Namespace string                  `yaml:"namespace"`
	Name      string                  `yaml:"name"`
	Key       string                  `yaml:"key"`
}

type SourceK8sContextConfigT struct {
	InCluster      bool   `yaml:"inCluster"`
	ConfigFilepath string `yaml:"configFilepath"`
	MasterUrl      string `yaml:"masterUrl"`
}

//--------------------------------------------------------------
// TARGET CONFIG
//--------------------------------------------------------------

type TargetT struct {
	Path string `yaml:"path"`
	File string `yaml:"file"`
	Mode uint32 `yaml:"mode"`
}

//--------------------------------------------------------------
// BEHAVIOR CONFIG
//--------------------------------------------------------------

type BehaviorConfigT struct {
	Conditions []ConditionConfigT `yaml:"conditions,omitempty"`
	Actions    []ActionConfigT    `yaml:"actions,omitempty"`
}

type ConditionConfigT struct {
	Name      string `yaml:"name"`
	Mandatory bool   `yaml:"mandatory"`
	Template  string `yaml:"template"`
	Expect    string `yaml:"expect"`
}

type ActionConfigT struct {
	Name string           `yaml:"name"`
	On   string           `yaml:"on"`
	In   string           `yaml:"in"`
	Cmd  []string         `yaml:"cmd"`
	K8s  ActionK8sConfigT `yaml:"k8s"`
}

type ActionK8sConfigT struct {
	Context   SourceK8sContextConfigT `yaml:"context"`
	Namespace string                  `yaml:"namespace"`
	Pod       string                  `yaml:"pod"`
	Container string                  `yaml:"container"`
	Command   []string                `yaml:"command"`
}
