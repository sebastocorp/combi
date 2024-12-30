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
	c.log.Info("init combi", map[string]any{})

	var err error
	for {
		c.log.Debug("waiting sync", map[string]any{})
		time.Sleep(c.syncTime)

		updatedList := []bool{}
		var updated bool
		for _, sv := range c.srcs {
			updated, err = sv.SyncConfig()
			if err != nil {
				break
			}
			updatedList = append(updatedList, updated)
		}

		if !slices.Contains(updatedList, true) {
			c.log.Debug("no updates in sources", map[string]any{})
			continue
		}
	}
}

func (c *CombiT) Stop() {
	c.log.Info("stop combi", map[string]any{})
}
