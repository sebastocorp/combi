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
	fileBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return c, err
	}

	conf, err := config.ParseConfig(fileBytes)
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
	c.log.Info("init combi", map[string]any{})

	var err error
	for {
		extraLogFields := map[string]any{}
		c.log.Debug("waiting next sync", extraLogFields)

		time.Sleep(c.syncTime)
		c.log.Info("init config sync", extraLogFields)

		// sources sync
		updatedList := []bool{}
		for _, sv := range c.srcs {
			globals.SetLogField(extraLogFields, globals.LogKeySourceName, sv.GetName())

			var updated bool
			updated, err = sv.SyncConfig()
			if err != nil {
				globals.SetLogField(extraLogFields, globals.LogKeyError, err.Error())
				c.log.Error("source sync failed", extraLogFields)
				break
			}

			updatedList = append(updatedList, updated)
		}
		globals.RemoveLogField(extraLogFields, globals.LogKeySourceName)
		if err != nil {
			continue
		}

		tFileExist := true
		if _, err = os.Stat(c.target.filepath); err != nil {
			if !os.IsNotExist(err) {
				globals.SetLogField(extraLogFields, globals.LogKeyError, err.Error())
				c.log.Error("unable to check target file", extraLogFields)
				continue
			}
			tFileExist = false
		}

		if !slices.Contains(updatedList, true) && tFileExist {
			c.log.Debug("no updates in sources", extraLogFields)
			continue
		}

		// sources decode and merge
		cfgResult := map[string]any{}
		for _, sv := range c.srcs {
			globals.SetLogField(extraLogFields, globals.LogKeySourceName, sv.GetName())

			var cfgBytes []byte
			cfgBytes, err = sv.GetConfig()
			if err != nil {
				globals.SetLogField(extraLogFields, globals.LogKeyError, err.Error())
				c.log.Error("unable to get source", extraLogFields)
				break
			}

			var cfg map[string]any
			cfg, err = c.encoder.DecodeConfigBytes(cfgBytes)
			if err != nil {
				globals.SetLogField(extraLogFields, globals.LogKeyError, err.Error())
				c.log.Error("unable to decode source", extraLogFields)
				break
			}

			c.encoder.MergeConfigs(cfgResult, cfg)
		}
		if err != nil {
			continue
		}
		globals.RemoveLogField(extraLogFields, globals.LogKeySourceName)

		// config conditions check
		condsResult := config.ConfigOnValueSUCCESS
		for _, cv := range c.conds {
			globals.RemoveLogField(extraLogFields, globals.LogKeyConditionResult)
			globals.SetLogField(extraLogFields, globals.LogKeyCondition, cv)

			var succ bool
			succ, err = cv.Eval(cfgResult)
			if err != nil {
				globals.SetLogField(extraLogFields, globals.LogKeyError, err.Error())
				c.log.Error("unable to evaluate condition", extraLogFields)
				break
			}

			globals.SetLogField(extraLogFields, globals.LogKeyConditionResult, globals.LogValueConditionResultSUCCESS)
			if !succ {
				globals.SetLogField(extraLogFields, globals.LogKeyConditionResult, globals.LogValueConditionResultFAIL)
				if cv.Mandatory {
					condsResult = config.ConfigOnValueFAILURE
				}
			}
			c.log.Debug("condition evaluated", extraLogFields)
		}
		if err != nil {
			continue
		}
		globals.RemoveLogField(extraLogFields, globals.LogKeyCondition)
		globals.RemoveLogField(extraLogFields, globals.LogKeyConditionResult)

		// config encode and create
		cfgResiltStr := c.encoder.EncodeConfigString(cfgResult)
		err = os.WriteFile(c.target.filepath, []byte(cfgResiltStr), c.target.mode)
		if err != nil {
			globals.SetLogField(extraLogFields, globals.LogKeyError, err.Error())
			c.log.Error("unable to create target file", extraLogFields)
			continue
		}

		// execute actions

		for _, av := range c.acts {
			globals.SetLogField(extraLogFields, globals.LogKeyAction, av)

			if av.On == condsResult {
				err = av.Exec()
				if err != nil {
					globals.SetLogField(extraLogFields, globals.LogKeyError, err.Error())
					c.log.Error("unable to execute action", extraLogFields)
					break
				}
				c.log.Debug("action executed", extraLogFields)
			}
		}
		if err != nil {
			continue
		}
		globals.RemoveLogField(extraLogFields, globals.LogKeyAction)

		c.log.Info("success in config sync", extraLogFields)
	}
}

func (c *CombiT) Stop() {
	c.log.Info("stop combi", map[string]any{})
}
