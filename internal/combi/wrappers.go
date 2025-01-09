package combi

import (
	"combi/api/v1alpha3"
	"combi/internal/encoding"
	"combi/internal/logger"
	"combi/internal/sources"
	"combi/internal/utils"
	"io/fs"
	"os"
	"path/filepath"
)

// setup TODO
func (c *CombiT) setup(conf *v1alpha3.CombiConfigT) error {
	c.log = logger.NewLogger(logger.GetLevel(conf.Logger.Level))
	c.syncTime = conf.Behavior.SyncTime
	c.encoder = encoding.GetEncoder(conf.Kind)

	// Target setup

	err := os.MkdirAll(conf.Behavior.Target.Path, fs.FileMode(conf.Behavior.Target.Mode))
	if err != nil {
		return err
	}

	c.target.filepath = filepath.Join(conf.Behavior.Target.Path, conf.Behavior.Target.File)
	c.target.mode = fs.FileMode(conf.Behavior.Target.Mode)

	// Sources setup

	for _, sv := range conf.Sources {
		var hashKey string
		hashKey, err = utils.GenHashString(sv.Type, sv.Name)
		if err != nil {
			return err
		}

		srcpath := filepath.Join(conf.Behavior.TmpObjs.Path, hashKey)
		err = os.MkdirAll(srcpath, fs.FileMode(conf.Behavior.TmpObjs.Mode))
		if err != nil {
			return err
		}

		var src sources.SourceT
		src, err = sources.GetSource(sv, srcpath)
		if err != nil {
			return err
		}
		c.srcs = append(c.srcs, src)
	}

	for _, cv := range conf.Behavior.Conditions {
		cond := NewCondition(cv)
		c.conds = append(c.conds, cond)
	}

	for _, av := range conf.Behavior.Actions {
		act := NewAction(av)
		c.acts = append(c.acts, act)
	}

	return nil
}
