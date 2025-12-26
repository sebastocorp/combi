package combi

import (
	"io/fs"
	"os"
	"text/template"
	"time"

	"combi/internal/globals"
	"combi/internal/logger"
	"combi/internal/sets/actions"
	"combi/internal/sets/conditions"
	"combi/internal/sets/credentials"
	"combi/internal/sets/sources"
	"combi/internal/utils"
)

const (
	TargetBuildTypeSOURCE   = "SOURCE"
	TargetBuildTypeTEMPLATE = "TEMPLATE"
)

type CombiT struct {
	log      logger.LoggerT
	syncTime time.Duration

	creds  *credentials.SetT
	srcs   *sources.SetT
	target TargetT
}

type TargetT struct {
	encType string
	build   BuildT
	cons    *conditions.SetT
	acts    *actions.SetT
}

type BuildT struct {
	bType string
	src   string
	tmpl  *template.Template
	file  string
	mode  fs.FileMode
}

// NewCombi TODO
func NewCombi(configFilePath string) (c *CombiT, err error) {
	cfgBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return c, err
	}

	conf, err := parseConfig(cfgBytes)
	if err != nil {
		return c, err
	}

	c = &CombiT{}
	err = c.setup(conf)
	if err != nil {
		return c, err
	}

	return c, err
}

// Run TODO
func (c *CombiT) Run() {
	c.log.Info("init combi", nil)

	var err error
	for {
		var extraLogFields logger.ExtraFieldsT = globals.GetLogCommonFields()
		c.log.Debug("waiting next sync", extraLogFields)

		time.Sleep(c.syncTime)
		c.log.Info("init config sync", extraLogFields)

		// sync sources
		var updated bool
		updated, err = c.srcs.Sync()
		if err != nil {
			extraLogFields.Set(globals.LogKeyError, err.Error())
			c.log.Error("sources sync failed", extraLogFields)
			continue
		}
		if !updated {
			if updated = !utils.FileExists(c.target.build.file); !updated {
				c.log.Debug("no updates in sources", extraLogFields)
				continue
			}
		}

		// get result file
		var cfgResult configResultT
		switch c.target.build.bType {
		case TargetBuildTypeSOURCE:
			{
				cfgResult, err = c.getConfigFromSource()
				if err != nil {
					extraLogFields.Set(globals.LogKeyError, err.Error())
					c.log.Error("get config from source failed", extraLogFields)
					continue
				}
			}
		case TargetBuildTypeTEMPLATE:
			{
				cfgResult, err = c.getConfigFromTemplate()
				if err != nil {
					extraLogFields.Set(globals.LogKeyError, err.Error())
					c.log.Error("get config from template failed", extraLogFields)
					continue
				}
			}
		}

		// check config conditions
		c.log.Debug("evaluate condition set", extraLogFields)
		var csr conditions.ResultT
		csr, err = c.target.cons.Evaluate(cfgResult.Map)
		extraLogFields.Set(globals.LogKeyConditionSet, csr)
		if err != nil {
			extraLogFields.Set(globals.LogKeyError, err.Error())
			c.log.Error("unable to evaluate condition set", extraLogFields)
			extraLogFields.Del(globals.LogKeyError)
			extraLogFields.Del(globals.LogKeyConditionSet)
			continue
		}
		c.log.Debug("condition set evaluated", extraLogFields)
		extraLogFields.Del(globals.LogKeyConditionSet)

		// config encode and create target file
		if csr.Status == conditions.StatusSuccess {
			err = os.WriteFile(c.target.build.file, cfgResult.Data, c.target.build.mode)
			if err != nil {
				extraLogFields.Set(globals.LogKeyError, err.Error())
				c.log.Error("unable to create target file", extraLogFields)
				extraLogFields.Del(globals.LogKeyError)
				continue
			}
		}

		// execute actions
		c.log.Debug("executing action set", extraLogFields)
		var asr actions.ResultT
		asr, err = c.target.acts.Execute(csr.Status)
		extraLogFields.Set(globals.LogKeyActionSet, asr)
		if err != nil {
			extraLogFields.Set(globals.LogKeyError, err.Error())
			c.log.Error("unable to execute action set", extraLogFields)
			extraLogFields.Del(globals.LogKeyError)
			extraLogFields.Del(globals.LogKeyActionSet)
			continue
		}
		c.log.Debug("action set executed", extraLogFields)
		extraLogFields.Del(globals.LogKeyActionSet)

		c.log.Info("success in config sync", extraLogFields)
	}
}

// Stop TODO
func (c *CombiT) Stop() {
	c.log.Info("stop combi", map[string]any{})
}
