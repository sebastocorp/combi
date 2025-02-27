package conditionset

import (
	"combi/internal/tmpl"
	"regexp"
)

const (
	StatusSuccess = "success"
	StatusFail    = "fail"
)

type ConditionSetT struct {
	set []conditionT
}

type OptionsT struct {
	Name      string
	Mandatory bool
	Tmpl      string
	Expect    string
}

type ResultT struct {
	Status string             `json:"status"`
	Crs    []ConditionResultT `json:"conditions"`
}

func NewConditionSet() (cs *ConditionSetT, err error) {
	cs = &ConditionSetT{}
	return cs, err
}

func (cs *ConditionSetT) CreateAdd(ops OptionsT) (err error) {
	c := conditionT{
		Mandatory: ops.Mandatory,
	}

	c.expect, err = regexp.Compile(ops.Expect)
	if err != nil {
		return err
	}

	c.tmpl, err = tmpl.NewTemplate(ops.Name, ops.Tmpl)
	if err != nil {
		return err
	}

	cs.set = append(cs.set, c)
	return err
}

func (cs *ConditionSetT) Evaluate(srcs map[string]any) (r ResultT, err error) {
	r.Status = StatusSuccess
	for ci := range cs.set {
		var cr ConditionResultT
		cr, err = cs.set[ci].eval(srcs)
		r.Crs = append(r.Crs, cr)
		if err != nil {
			r.Status = StatusFail
			return r, err
		}

		if cr.Mandatory && !cr.Match {
			r.Status = StatusFail
		}
	}

	return r, err
}
