package combi

import (
	"bytes"
	"combi/api/v1alpha4"
	"combi/internal/tmpl"
	"regexp"
	"text/template"
)

type ConditionT struct {
	Name      string `json:"name"`
	Mandatory bool   `json:"mandatory"`

	tmpl   *template.Template `json:"-"`
	expect *regexp.Regexp     `json:"-"`
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

func (c *ConditionT) Eval(cfg map[string]any) (success bool, err error) {
	var result string
	result, err = c.evaluate(cfg)
	if err != nil {
		return success, err
	}

	success = c.expect.MatchString(result)

	return success, err
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
