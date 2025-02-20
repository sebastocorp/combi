package combi

import (
	"combi/api/v1alpha4"
	"combi/internal/encoding"
	"combi/internal/logger"
	"combi/internal/sources"
	"combi/internal/utils"
	"io/fs"
	"os"
	"path/filepath"
)

// setup TODO
func (c *CombiT) setup(cfg *v1alpha4.CombiConfigT) error {
	c.log = logger.NewLogger(logger.GetLevel(cfg.Settings.Logger.Level))
	c.syncTime = cfg.Settings.SyncTime
	c.encoder = encoding.GetEncoder(cfg.Kind)

	// Target setup

	err := os.MkdirAll(cfg.Settings.Target.Path, fs.FileMode(cfg.Settings.Target.Mode))
	if err != nil {
		return err
	}

	c.target.filepath = filepath.Join(cfg.Settings.Target.Path, cfg.Settings.Target.File)
	c.target.mode = fs.FileMode(cfg.Settings.Target.Mode)

	// Sources setup

	for si := range cfg.Sources {
		var hashKey string
		hashKey, err = utils.GenHashString(cfg.Sources[si].Type, cfg.Sources[si].Name)
		if err != nil {
			return err
		}

		srcpath := filepath.Join(cfg.Settings.TmpObjs.Path, hashKey)
		err = os.MkdirAll(srcpath, fs.FileMode(cfg.Settings.TmpObjs.Mode))
		if err != nil {
			return err
		}

		var src sources.SourceT
		src, err = sources.GetSource(cfg.Sources[si], srcpath)
		if err != nil {
			return err
		}
		c.srcs = append(c.srcs, src)
	}

	for ci := range cfg.Behavior.Conditions {
		var cond ConditionT
		cond, err = NewCondition(cfg.Behavior.Conditions[ci])
		if err != nil {
			return err
		}
		c.conds = append(c.conds, cond)
	}

	for ai := range cfg.Behavior.Actions {
		var act ActionT
		act, err = NewAction(cfg.Behavior.Actions[ai])
		if err != nil {
			return err
		}

		c.acts = append(c.acts, act)
	}

	return nil
}
