package combi

import (
	"combi/api/v1alpha3"
	"combi/internal/encoding"
	"combi/internal/globals"
	"combi/internal/logger"
	"combi/internal/sources"
	"combi/internal/utils"
	"io/fs"
	"os"
	"path/filepath"
)

// setup TODO
func (c *CombiT) setup(conf *v1alpha3.CombiConfigT) error {
	c.log = logger.NewLogger(logger.GetLevel(conf.Logger.Level), globals.GetLogCommonFields())
	c.syncTime = conf.Behavior.SyncTime
	c.encoder = encoding.GetEncoder(conf.Kind)

	// Target setup

	err := os.MkdirAll(conf.Behavior.Target.Path, fs.FileMode(conf.Behavior.Target.Mode))
	if err != nil {
		return err
	}

	c.targetFilepath = conf.Behavior.Target.Path + "/" + conf.Behavior.Target.File

	// Sources setup

	c.srcs = map[string]sources.SourceT{}
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
		c.srcs[hashKey] = src
	}

	return nil
}
