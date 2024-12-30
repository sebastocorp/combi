package v1alpha3

import "time"

type CombiConfigT struct {
	Kind     string          `yaml:"kind"`
	Logger   LoggerConfigT   `yaml:"logger"`
	Sources  []SourceConfigT `yaml:"sources"`
	Behavior BehaviorConfigT `yaml:"behavior"`
}

type LoggerConfigT struct {
	Level string `yaml:"level"`
}

//--------------------------------------------------------------
// SOURCE CONFIG
//--------------------------------------------------------------

type SourceConfigT struct {
	Name       string           `yaml:"name"`
	Type       string           `yaml:"type"` // values: raw|file|git|k8s
	Raw        string           `yaml:"raw,omitempty"`
	File       string           `yaml:"file,omitempty"`
	Git        SourceGitConfigT `yaml:"git,omitempty"`
	Kubernetes SourceK8sConfigT `yaml:"k8s,omitempty"`
}

type SourceGitConfigT struct {
	SshUrl         string `yaml:"sshUrl"`
	SshKeyFilepath string `yaml:"sshKeyFilepath"`
	Branch         string `yaml:"branch"`
	Filepath       string `yaml:"filepath"`
}

type SourceK8sConfigT struct {
	Kind      string `yaml:"kind"`
	Namespace string `yaml:"namespace"`
	Name      string `yaml:"name"`
	Key       string `yaml:"key"`
}

//--------------------------------------------------------------
// BEHAVIOR CONFIG
//--------------------------------------------------------------

type BehaviorConfigT struct {
	SyncTime   time.Duration      `yaml:"syncTime"`
	Target     TargetConfigT      `yaml:"target"`
	TmpObjs    TmpObjsConfigT     `yaml:"tmpObjs"`
	Conditions []ConditionConfigT `yaml:"conditions,omitempty"`
	Actions    []ActionConfigT    `yaml:"actions,omitempty"`
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

type ConditionConfigT struct {
	Name      string `yaml:"name"`
	Mandatory bool   `yaml:"mandatory"`
	Template  string `yaml:"template"`
	Value     string `yaml:"value"`
}

type ActionConfigT struct {
	Name            string   `yaml:"name"`
	ConditionResult string   `yaml:"conditionResult"`
	Command         []string `yaml:"command"`
	Script          string   `yaml:"script"`
}
