package combi

import (
	"combi/api/v1alpha3"
	"combi/internal/template"
)

type ConditionT struct {
	Name      string `json:"name"`
	Mandatory bool   `json:"mandatory"`

	tmpl string `json:"-"`
	val  string `json:"-"`
}

func NewCondition(condition v1alpha3.ConditionConfigT) ConditionT {
	return ConditionT{
		Name:      condition.Name,
		Mandatory: condition.Mandatory,

		tmpl: condition.Template,
		val:  condition.Value,
	}
}

func (c *ConditionT) Eval(cfg map[string]any) (success bool, err error) {
	var result string
	result, err = template.EvaluateTemplate(c.tmpl, cfg)
	if err != nil {
		return success, err
	}

	success = result == c.val

	return success, err
}
