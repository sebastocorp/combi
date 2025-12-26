package actions

import (
	"fmt"
)

const (
	TypeK8S   = "K8S"
	TypeLOCAL = "LOCAL"
)

type SetT struct {
	as []ActionI
}

type ActionI interface {
	getOn() string
	exec() (r ActionResultT, err error)
}

type OptionsT struct {
	Name    string
	Type    string
	On      string
	CredRef any

	K8s OptionsK8sT
	Cmd OptionsCmdT
}

type ResultT struct {
	Ars []ActionResultT `json:"actions"`
}

type ActionResultT struct {
	Name   string `json:"name"`
	On     string `json:"on"`
	Cmd    string `json:"cmd"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func NewSet() (s *SetT, err error) {
	s = &SetT{}
	return s, err
}

func (s *SetT) Add(ops OptionsT) (err error) {
	switch ops.Type {
	case TypeK8S:
		{
			var act *K8sActionT
			act, err = newK8sAction(ops)
			if err != nil {
				return err
			}
			s.as = append(s.as, act)
		}
	case TypeLOCAL:
		{
			var act *CmdActionT
			act, err = newCmdAction(ops)
			if err != nil {
				return err
			}
			s.as = append(s.as, act)
		}
	default:
		{
			err = fmt.Errorf("unsupported '%s' action type", ops.Type)
			return err
		}
	}

	return err
}

func (s *SetT) Execute(on string) (r ResultT, err error) {
	for acti := range s.as {
		if s.as[acti].getOn() == on {
			var ar ActionResultT
			ar, err = s.as[acti].exec()
			r.Ars = append(r.Ars, ar)
			if err != nil {
				return r, err
			}
		}
	}

	return r, err
}
