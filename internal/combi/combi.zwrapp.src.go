package combi

import (
	"bytes"
	"combi/internal/encoders"
	"combi/internal/sets/sources"
)

type configResultT struct {
	Data []byte
	Map  map[string]any
}

func (c *CombiT) getConfigFromSource() (result configResultT, err error) {
	srcd, err := c.srcs.GetByName(c.target.build.src)
	if err != nil {
		return result, err
	}
	result.Data = srcd.Data

	result.Map, err = encoders.Encoders[c.target.encType].Decode(result.Data)
	return result, err
}

func (c *CombiT) getConfigFromTemplate() (result configResultT, err error) {
	srcMaps := make(map[string]any)
	for si := range c.srcs.Length() {
		var srcd sources.SourceDataT
		srcd, err = c.srcs.GetByIndex(si)
		if err != nil {

			return result, err
		}

		var cfg map[string]any
		cfg, err = encoders.Encoders[srcd.EncType].Decode(srcd.Data)
		if err != nil {
			return result, err
		}

		srcMaps[srcd.Name] = cfg
	}

	buffer := new(bytes.Buffer)
	err = c.target.build.tmpl.Execute(buffer, srcMaps)
	if err != nil {
		return result, err
	}

	result.Data = buffer.Bytes()
	result.Map, err = encoders.Encoders[c.target.encType].Decode(result.Data)

	return result, err
}
