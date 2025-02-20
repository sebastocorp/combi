package v1alpha4

import "time"

type CombiConfigT struct {
	Kind     string          `yaml:"kind"`
	Settings SettingsConfigT `yaml:"settings"`
	Sources  []SourceConfigT `yaml:"sources"`
	Behavior BehaviorConfigT `yaml:"behavior"`
}

//--------------------------------------------------------------
// SETTINGS CONFIG
//--------------------------------------------------------------

type SettingsConfigT struct {
	Logger   LoggerConfigT  `yaml:"logger"`
	SyncTime time.Duration  `yaml:"syncTime"`
	Target   TargetConfigT  `yaml:"target"`
	TmpObjs  TmpObjsConfigT `yaml:"tmpObjs"`
}

type LoggerConfigT struct {
	Level string `yaml:"level"`
}

type TmpObjsConfigT struct {
	Path string `yaml:"path"`
	Mode uint32 `yaml:"mode"`
}

type TargetConfigT struct {
	Path string `yaml:"path"`
	File string `yaml:"file"`
	Mode uint32 `yaml:"mode"`
}

//--------------------------------------------------------------
// SOURCE CONFIG
//--------------------------------------------------------------

type SourceConfigT struct {
	Name string           `yaml:"name"`
	Type string           `yaml:"type"` // values: RAW|FILE|GIT|K8S
	Raw  string           `yaml:"raw,omitempty"`
	File string           `yaml:"file,omitempty"`
	Git  SourceGitConfigT `yaml:"git,omitempty"`
	K8s  SourceK8sConfigT `yaml:"k8s,omitempty"`
}

type SourceGitConfigT struct {
	SshUrl         string `yaml:"sshUrl"`
	SshKeyFilepath string `yaml:"sshKeyFilepath"`
	Branch         string `yaml:"branch"`
	Filepath       string `yaml:"filepath"`
}

type SourceK8sConfigT struct {
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
	Name    string           `yaml:"name"`
	On      string           `yaml:"on"`
	Command []string         `yaml:"command"`
	K8s     ActionK8sConfigT `yaml:"k8s"`
}

type ActionK8sConfigT struct {
	Context   SourceK8sContextConfigT `yaml:"context"`
	Namespace string                  `yaml:"namespace"`
	Pod       string                  `yaml:"pod"`
	Container string                  `yaml:"container"`
	Command   []string                `yaml:"command"`
}
