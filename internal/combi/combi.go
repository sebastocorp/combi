package combi

import (
	"io/fs"
	"os"
	"slices"
	"time"

	"combi/internal/config"
	"combi/internal/encoding"
	"combi/internal/globals"
	"combi/internal/logger"
	"combi/internal/sources"
)

type CombiT struct {
	log logger.LoggerT

	syncTime time.Duration
	target   TargetT
	encoder  encoding.EncoderT
	srcs     []sources.SourceT
	conds    []ConditionT
	acts     []ActionT
}

type TargetT struct {
	filepath string
	mode     fs.FileMode
}

func NewCombi(configFilePath string) (c *CombiT, err error) {
	cfgBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return c, err
	}

	conf, err := config.ParseConfig(cfgBytes)
	if err != nil {
		return c, err
	}

	c = &CombiT{}
	err = c.setup(&conf)
	if err != nil {
		return c, err
	}

	return c, err
}

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
		for _, sv := range c.srcs {
			extraLogFields.Set(globals.LogKeySourceName, sv.GetName())

			var cfgBytes []byte
			cfgBytes, err = sv.GetConfig()
			if err != nil {
				extraLogFields.Set(globals.LogKeyError, err.Error())
				c.log.Error("unable to get source", extraLogFields)
				break
			}

			var cfg map[string]any
			cfg, err = c.encoder.DecodeConfig(cfgBytes)
			if err != nil {
				extraLogFields.Set(globals.LogKeyError, err.Error())
				c.log.Error("unable to decode source", extraLogFields)
				break
			}

			c.encoder.MergeConfigs(cfgResult, cfg)
		}
		if err != nil {
			continue
		}
		extraLogFields.Del(globals.LogKeySourceName)

		// check config conditions
		condsResult := config.ConfigOnValueSUCCESS
		for _, cv := range c.conds {
			extraLogFields.Del(globals.LogKeyConditionResult)
			extraLogFields.Set(globals.LogKeyCondition, cv)

			var succ bool
			succ, err = cv.Eval(cfgResult)
			if err != nil {
				extraLogFields.Set(globals.LogKeyError, err.Error())
				c.log.Error("unable to evaluate condition", extraLogFields)
				break
			}

			extraLogFields.Set(globals.LogKeyConditionResult, globals.LogValueConditionResultSUCCESS)
			if !succ {
				extraLogFields.Set(globals.LogKeyConditionResult, globals.LogValueConditionResultFAIL)
				if cv.Mandatory {
					condsResult = config.ConfigOnValueFAILURE
				}
			}
			c.log.Debug("condition evaluated", extraLogFields)
		}
		if err != nil {
			continue
		}
		extraLogFields.Del(globals.LogKeyCondition)
		extraLogFields.Del(globals.LogKeyConditionResult)

		// config encode and create target file
		if condsResult == config.ConfigOnValueSUCCESS {
			var cfgResultBytes []byte
			cfgResultBytes, err = c.encoder.EncodeConfig(cfgResult)
			if err != nil {
				extraLogFields.Set(globals.LogKeyError, err.Error())
				c.log.Error("unable to generate config", extraLogFields)
				continue
			}

			err = os.WriteFile(c.target.filepath, cfgResultBytes, c.target.mode)
			if err != nil {
				extraLogFields.Set(globals.LogKeyError, err.Error())
				c.log.Error("unable to create target file", extraLogFields)
				continue
			}
		}

		// execute actions
		for _, av := range c.acts {
			extraLogFields.Set(globals.LogKeyAction, av)

			if av.On == condsResult {
				var outBytes, errBytes []byte
				outBytes, errBytes, err = av.Exec()
				if err != nil {
					extraLogFields.Set(globals.LogKeyError, err.Error())
					c.log.Error("unable to execute action", extraLogFields)
					break
				}
				extraLogFields.Set(globals.LogKeyActionStdout, string(outBytes))
				extraLogFields.Set(globals.LogKeyActionStderr, string(errBytes))
				c.log.Debug("action executed", extraLogFields)
				extraLogFields.Del(globals.LogKeyActionStdout)
				extraLogFields.Del(globals.LogKeyActionStderr)
			}
		}
		if err != nil {
			continue
		}
		extraLogFields.Del(globals.LogKeyAction)

		c.log.Info("success in config sync", extraLogFields)
	}
}

func (c *CombiT) Stop() {
	c.log.Info("stop combi", map[string]any{})
}
