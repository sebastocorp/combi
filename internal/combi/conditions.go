package combi

import (
	"bytes"
	"combi/api/v1alpha4"
	"combi/internal/tmpl"
	"regexp"
	"text/template"
)

const (
	ConditionStatusSuccess = "success"
	ConditionStatusFail    = "fail"
)

type ConditionT struct {
	Name      string `json:"name"`
	Mandatory bool   `json:"mandatory"`

	tmpl   *template.Template `json:"-"`
	expect *regexp.Regexp     `json:"-"`
}

type ConditionResultT struct {
	Status     string `json:"status"`
	TmplResult string `json:"tmplResult"`
	Expect     string `json:"expect"`
}

func NewCondition(condCfg v1alpha4.ConditionConfigT) (cond ConditionT, err error) {
	cond = ConditionT{
		Name:      condCfg.Name,
		Mandatory: condCfg.Mandatory,

		expect: regexp.MustCompile(condCfg.Expect),
	}

	cond.tmpl, err = tmpl.NewTemplate(condCfg.Name, condCfg.Template)
	return cond, err
}

func (c *ConditionT) Eval(cfg map[string]any) (result ConditionResultT, err error) {
	result.Status = ConditionStatusFail
	result.Expect = c.expect.String()
	result.TmplResult, err = c.evaluate(cfg)
	if err != nil {
		return result, err
	}

	if c.expect.MatchString(result.TmplResult) {
		result.Status = ConditionStatusSuccess
	}

	return result, err
}

func (c *ConditionT) evaluate(srcs map[string]any) (result string, err error) {
	// Create a new buffer to store the templating result
	buffer := new(bytes.Buffer)
	err = c.tmpl.Execute(buffer, srcs)
	if err != nil {
		return result, err
	}

	return buffer.String(), nil
}
