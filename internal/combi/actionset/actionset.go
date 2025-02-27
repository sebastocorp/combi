package actionset

import (
	"reflect"

	"combi/internal/utils"
)

type ActionSetT struct {
	set []actionT
}

type OptionsT struct {
	Name string `json:"name"`
	On   string `json:"on"`

	K8s OptionsK8sT `json:"-"`
	Cmd []string    `json:"cmd"`
}

type OptionsK8sT struct {
	InCluster      bool
	ConfigFilepath string
	MasterUrl      string

	Namespace string
	Pod       string
	Container string
}

type ResultT struct {
	Ars []ActionResultT `json:"actions"`
}

func NewActionSet() (as *ActionSetT, err error) {
	as = &ActionSetT{}
	return as, err
}

func (as *ActionSetT) CreateAdd(ops OptionsT) (err error) {
	a := actionT{
		Name: ops.Name,
		On:   ops.On,
		cmd:  ops.Cmd,
	}

	if reflect.TypeOf(ops.K8s) != reflect.TypeOf(OptionsK8sT{}) {
		a.k8s = actionK8sT{
			namespace: ops.K8s.Namespace,
			pod:       ops.K8s.Pod,
			container: ops.K8s.Container,
		}

		a.k8s.cfg, err = utils.GetK8sConfig(ops.K8s.InCluster, ops.K8s.ConfigFilepath, ops.K8s.MasterUrl)
		if err != nil {
			return err
		}

		a.k8s.client, err = utils.GetK8sClient(a.k8s.cfg)
		if err != nil {
			return err
		}
	}

	as.set = append(as.set, a)

	return err
}

func (as *ActionSetT) Execute(on string) (r ResultT, err error) {
	for ai := range as.set {
		if as.set[ai].On == on {
			var ar ActionResultT
			ar, err = as.set[ai].exec()
			r.Ars = append(r.Ars, ar)
			if err != nil {
				return r, err
			}
		}
	}

	return r, err
}
