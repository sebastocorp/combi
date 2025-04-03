package combi

import (
	"io/fs"
	"os"
	"time"

	"combi/internal/encoding"
	"combi/internal/globals"
	"combi/internal/logger"
	"combi/internal/sets/actions"
	"combi/internal/sets/conditions"
	"combi/internal/sets/credentials"
	"combi/internal/sets/sources"
)

const (
	CombiDirMode   fs.FileMode = 0644
	CombiFilesMode fs.FileMode = 0755
)

type CombiT struct {
	log logger.LoggerT

	syncTime time.Duration
	target   TargetT
	encoder  encoding.EncoderT

	creds *credentials.SetT
	srcs  *sources.SetT
	cons  *conditions.SetT
	acts  *actions.SetT
}

type TargetT struct {
	filepath string
	mode     fs.FileMode
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

	os.Exit(0)

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

		tFileExist := true
		if _, err = os.Stat(c.target.filepath); err != nil {
			if !os.IsNotExist(err) {
				extraLogFields.Set(globals.LogKeyError, err.Error())
				c.log.Error("unable to check target file", extraLogFields)
				extraLogFields.Del(globals.LogKeyError)
				continue
			}
			tFileExist = false
		}

		if !updated && tFileExist {
			c.log.Debug("no updates in sources", extraLogFields)
			continue
		}

		// decode and merge sources
		cfgResult := map[string]any{}
		cfgSrcBytes := []byte{}
		for si := range c.srcs.Length() {
			var cfgBytes []byte
			cfgBytes, err = c.srcs.GetByIndex(si)
			if err != nil {
				extraLogFields.Set(globals.LogKeyError, err.Error())
				c.log.Error("unable to get source", extraLogFields)
				extraLogFields.Del(globals.LogKeyError)
				break
			}

			var cfg map[string]any
			cfg, err = c.encoder.DecodeConfig(cfgBytes)
			if err != nil {
				extraLogFields.Set(globals.LogKeyError, err.Error())
				c.log.Error("unable to decode source", extraLogFields)
				extraLogFields.Del(globals.LogKeyError)
				break
			}

			c.encoder.MergeConfigs(cfgResult, cfg)

			if c.srcs.Length() == 1 {
				cfgSrcBytes = cfgBytes
			}
		}
		if err != nil {
			continue
		}
		extraLogFields.Del(globals.LogKeySourceName)

		// check config conditions
		c.log.Debug("evaluate condition set", extraLogFields)
		var csr conditions.ResultT
		csr, err = c.cons.Evaluate(cfgResult)
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
			var cfgResultBytes []byte
			if c.srcs.Length() != 1 {
				cfgResultBytes, err = c.encoder.EncodeConfig(cfgResult)
				if err != nil {
					extraLogFields.Set(globals.LogKeyError, err.Error())
					c.log.Error("unable to generate config", extraLogFields)
					extraLogFields.Del(globals.LogKeyError)
					continue
				}
			} else {
				cfgResultBytes = cfgSrcBytes
			}

			err = os.WriteFile(c.target.filepath, cfgResultBytes, c.target.mode)
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
		asr, err = c.acts.Execute(csr.Status)
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
