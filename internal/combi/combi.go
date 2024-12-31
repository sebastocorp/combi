package combi

import (
	"fmt"
	"os"
	"slices"
	"time"

	"combi/internal/config"
	"combi/internal/encoding"
	"combi/internal/logger"
	"combi/internal/sources"
)

type CombiT struct {
	log logger.LoggerT

	syncTime       time.Duration
	targetFilepath string
	encoder        encoding.EncoderT
	srcs           map[string]sources.SourceT
}

func NewCombi(configFilePath string) (c *CombiT, err error) {
	fileBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return c, err
	}

	conf, err := config.ParseConfig(fileBytes)
	if err != nil {
		return c, fmt.Errorf("unable to parse config file: %s", err.Error())
	}

	c = &CombiT{}
	err = c.setup(&conf)
	if err != nil {
		return c, err
	}

	//--------------------------------------------------------------
	//
	//--------------------------------------------------------------

	return c, err
}

func (c *CombiT) Run() {
	extraLogFile := map[string]any{}
	c.log.Info("init combi", extraLogFile)
	var err error
	for {
		c.log.Debug("waiting sync", extraLogFile)
		time.Sleep(c.syncTime)

		c.log.Info("init sources sync", extraLogFile)
		updatedList := []bool{}
		var updated bool
		for _, sv := range c.srcs {
			updated, err = sv.SyncConfig()
			if err != nil {
				c.log.Error("source sync failed", extraLogFile)
				break
			}
			updatedList = append(updatedList, updated)
		}

		if !slices.Contains(updatedList, true) {
			c.log.Debug("no updates in sources", extraLogFile)
			continue
		}
	}
}

func (c *CombiT) Stop() {
	c.log.Info("stop combi", map[string]any{})
}
