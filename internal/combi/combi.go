package combi

import (
	"io/fs"
	"os"
	"slices"
	"time"

	"combi/internal/credentials"
	"combi/internal/encoding"
	"combi/internal/globals"
	"combi/internal/logger"
	"combi/internal/sources"
	"combi/internal/target/actionset"
	"combi/internal/target/conditionset"
)

type CombiT struct {
	log logger.LoggerT

	syncTime time.Duration
	target   TargetT
	encoder  encoding.EncoderT
	srcs     []sources.SourceT
	cs       *conditionset.ConditionSetT
	as       *actionset.ActionSetT

	// new
	creds *credentials.CredentialSetT
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
		updatedList := []bool{}
		for _, sv := range c.srcs {
			extraLogFields.Set(globals.LogKeySourceName, sv.GetName())

			var updated bool
			updated, err = sv.SyncConfig()
			if err != nil {
				extraLogFields.Set(globals.LogKeyError, err.Error())
				c.log.Error("source sync failed", extraLogFields)
				break
			}

			updatedList = append(updatedList, updated)
		}
		extraLogFields.Del(globals.LogKeySourceName)
		if err != nil {
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

		if !slices.Contains(updatedList, true) && tFileExist {
			c.log.Debug("no updates in sources", extraLogFields)
			continue
		}

		// decode and merge sources
		cfgResult := map[string]any{}
		cfgSrcBytes := []byte{}
		for _, sv := range c.srcs {
			extraLogFields.Set(globals.LogKeySourceName, sv.GetName())

			var cfgBytes []byte
			cfgBytes, err = sv.GetConfig()
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

			if len(c.srcs) == 1 {
				cfgSrcBytes = cfgBytes
			}
		}
		if err != nil {
			continue
		}
		extraLogFields.Del(globals.LogKeySourceName)

		// check config conditions
		c.log.Debug("evaluate condition set", extraLogFields)
		var csr conditionset.ResultT
		csr, err = c.cs.Evaluate(cfgResult)
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
		if csr.Status == conditionset.StatusSuccess {
			var cfgResultBytes []byte
			if len(c.srcs) != 1 {
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
		var asr actionset.ResultT
		asr, err = c.as.Execute(csr.Status)
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
