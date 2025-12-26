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
	SyncTime   time.Duration `yaml:"syncTime"`
	WorkingDir string        `yaml:"workingDir"`
	Logger     LoggerT       `yaml:"logger"`

	Credentials []CredentialT `yaml:"credentials"`
	Sources     []SourceT     `yaml:"sources"`
	Target      TargetT       `yaml:"target"`
}

type LoggerT struct {
	Level string `yaml:"level"`
}

//--------------------------------------------------------------
// CREDENTIALS CONFIG
//--------------------------------------------------------------

type CredentialT struct {
	Name   string            `yaml:"name"`
	Type   string            `yaml:"type"`
	SshKey CredentialSshKeyT `yaml:"sshKey"`
	K8s    CredentialK8sT    `yaml:"k8s"`
}

type CredentialSshKeyT struct {
	User       string `yaml:"user"`
	SshKeyFile string `yaml:"sshKeyFile"`
	Password   string `yaml:"password"`
}

type CredentialK8sT struct {
	InCluster      bool   `yaml:"inCluster"`
	KubeconfigPath string `yaml:"kubeconfigPath"`
	MasterUrl      string `yaml:"masterUrl"`
}

//--------------------------------------------------------------
// SOURCES CONFIG
//--------------------------------------------------------------

type SourceT struct {
	Name       string     `yaml:"name"`
	Type       string     `yaml:"type"`    // values: FILE_RAW|FILE|GIT|K8S
	Encoder    string     `yaml:"encoder"` // values: YAML|JSON|NGINX|LIBCONFIG
	Credential string     `yaml:"credential,omitempty"`
	File       string     `yaml:"file,omitempty"`
	Git        SourceGitT `yaml:"git,omitempty"`
	K8s        SourceK8sT `yaml:"k8s,omitempty"`
}

type SourceGitT struct {
	SshUrl string `yaml:"sshUrl"`
	Branch string `yaml:"branch"`
	File   string `yaml:"file"`
}

type SourceK8sT struct {
	Kind      string `yaml:"kind"`
	Namespace string `yaml:"namespace"`
	Name      string `yaml:"name"`
	Key       string `yaml:"key"`
}

//--------------------------------------------------------------
// TARGET CONFIG
//--------------------------------------------------------------

type TargetT struct {
	Encoder    string             `yaml:"encoder"` // values: YAML|JSON|NGINX|LIBCONFIG
	Build      TargetBuildT       `yaml:"build"`
	Conditions []TargetConditionT `yaml:"conditions,omitempty"`
	Actions    []TargetActionT    `yaml:"actions,omitempty"`
}

type TargetBuildT struct {
	File     string `yaml:"file"`
	Mode     uint32 `yaml:"mode"`
	Type     string `yaml:"type"` // values: TEMPLATE|SOURCE
	Source   string `yaml:"source"`
	Template string `yaml:"template"`
}

type TargetConditionT struct {
	Name      string `yaml:"name"`
	Mandatory bool   `yaml:"mandatory"`
	Template  string `yaml:"template"`
	Expect    string `yaml:"expect"`
}

type TargetActionT struct {
	Name       string           `yaml:"name"`
	On         string           `yaml:"on"`   // values: SUCCESS|FAILURE
	Type       string           `yaml:"type"` // values: LOCAL|K8S
	Credential string           `yaml:"credential,omitempty"`
	Cmd        []string         `yaml:"cmd,omitempty"`
	K8s        TargetActionK8sT `yaml:"k8s,omitempty"`
}

type TargetActionK8sT struct {
	Credential string   `yaml:"credential"`
	Namespace  string   `yaml:"namespace"`
	Pod        string   `yaml:"pod"`
	Container  string   `yaml:"container"`
	Cmd        []string `yaml:"cmd"`
}
