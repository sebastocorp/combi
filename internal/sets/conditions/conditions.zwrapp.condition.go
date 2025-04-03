package conditions

import (
	"bytes"
	"regexp"
	"text/template"
)

type conditionT struct {
	Mandatory bool

	tmpl   *template.Template
	expect *regexp.Regexp
}

type ConditionResultT struct {
	Name      string `json:"name"`
	Mandatory bool   `json:"mandatory"`
	Match     bool   `json:"match"`
	Tmpl      string `json:"tmpl"`
	Expect    string `json:"expect"`
}

func (c *conditionT) eval(cfg map[string]any) (r ConditionResultT, err error) {
	r.Name = c.tmpl.Name()
	r.Mandatory = c.Mandatory
	r.Expect = c.expect.String()

	r.Tmpl, err = c.evaluate(cfg)
	if err != nil {
		return r, err
	}

	r.Match = c.expect.MatchString(r.Tmpl)

	return r, err
}

func (c *conditionT) evaluate(srcs map[string]any) (result string, err error) {
	// Create a new buffer to store the templating result
	buffer := new(bytes.Buffer)
	err = c.tmpl.Execute(buffer, srcs)
	if err != nil {
		return result, err
	}

	return buffer.String(), nil
}
