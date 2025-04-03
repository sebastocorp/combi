package conditions

import (
	"combi/internal/tmpl"
	"regexp"
)

const (
	StatusSuccess = "SUCCESS"
	StatusFail    = "FAIL"
)

type SetT struct {
	cs []conditionT
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

func NewSet() (s *SetT, err error) {
	s = &SetT{}
	return s, err
}

func (s *SetT) Add(ops OptionsT) (err error) {
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

	s.cs = append(s.cs, c)
	return err
}

func (s *SetT) Evaluate(srcs map[string]any) (r ResultT, err error) {
	r.Status = StatusSuccess
	for ci := range s.cs {
		var cr ConditionResultT
		cr, err = s.cs[ci].eval(srcs)
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
