package nginx

import (
	"fmt"
	"regexp"
	"strings"
)

type NginxT struct {
	ConfigStruct BlockContentT
}

// ----------------------------------------------------------------
// NGINX data structure
// ----------------------------------------------------------------

type BlockContentT struct {
	Directives []DirectiveT
	Blocks     []BlockT
}

type DirectiveT struct {
	Name  string
	Param string
	Value string
}

type BlockT struct {
	Name         string
	Params       string
	BlockContent BlockContentT
}

// ----------------------------------------------------------------
// Decode/Encode NGINX data structure
// ----------------------------------------------------------------

// Decode functions
func (e *NginxT) DecodeConfigBytes(configBytes []byte) (cfg map[string]any, err error) {
	// Remove one line comments in file
	// configStr := regexp.MustCompile(`#[^\n]*?\n`).ReplaceAllString(string(configBytes), "\n")
	configStr := regexp.MustCompile(`\n\s+`).ReplaceAllString(string(configBytes), "\n")

	// Format configuration to parse nginx config format by line
	configStr = strings.Join(strings.Fields(configStr), " ")
	configStr = strings.ReplaceAll(configStr, "\n", "")
	configStr = strings.ReplaceAll(configStr, "{", " {\n")
	configStr = strings.ReplaceAll(configStr, "}", "}\n")
	configStr = strings.ReplaceAll(configStr, ";", " ;\n")

	// Parse formatted nginx configuration
	configStrLines := strings.Split(configStr, "\n")
	configStruct := BlockContentT{}
	err = parseNginxBlockContent(&configStruct, configStrLines)
	cfg = e.ConfigToMap(&configStruct)
	return cfg, err
}

// func (e *NginxT) DecodeConfigBytes(configBytes []byte) (err error) {
// 	// Remove one line comments in file
// 	re := regexp.MustCompile(`#[^\n]*?\n`)
// 	configStr := re.ReplaceAllString(string(configBytes), "\n")

// 	// Format configuration to parse nginx config format by line
// 	configStr = strings.Join(strings.Fields(configStr), " ")
// 	configStr = strings.ReplaceAll(configStr, "\n", "")
// 	configStr = strings.ReplaceAll(configStr, "{", " {\n")
// 	configStr = strings.ReplaceAll(configStr, "}", "}\n")
// 	configStr = strings.ReplaceAll(configStr, ";", " ;\n")
// 	configStr = strings.ReplaceAll(configStr, "\n ", "\n")

// 	// Parse formatted nginx configuration
// 	configStrLines := strings.Split(configStr, "\n")
// 	err = parseNginxBlockContent(&e.ConfigStruct, configStrLines)
// 	return err
// }

func parseNginxBlockContent(blockContent *BlockContentT, blockContentLines []string) (err error) {
	// Parse block content
	for blockLineIndex := 0; blockLineIndex < len(blockContentLines); blockLineIndex++ {
		line := strings.TrimSpace(blockContentLines[blockLineIndex])
		// Skip empty strings (only to be carefull)
		if len(line) == 0 {
			continue
		}

		// Parse nginx directives
		if strings.HasSuffix(line, ";") {
			directiveParts := strings.Fields(line)
			subDirective := DirectiveT{
				Name:  directiveParts[0],
				Value: strings.Join(directiveParts[1:len(directiveParts)-1], " "),
			}
			if len(directiveParts) > 3 {
				subDirective.Param = directiveParts[1]
				subDirective.Value = strings.Join(directiveParts[2:len(directiveParts)-1], " ")
			}
			blockContent.Directives = append(blockContent.Directives, subDirective)
		}

		// Parse ngix blocks
		if strings.HasSuffix(line, "{") {
			startBlockLines := blockContentLines[blockLineIndex:]
			i, open := 0, 0
			for ; i < len(startBlockLines); i++ {
				if strings.HasSuffix(startBlockLines[i], "{") {
					open += 1
				}
				if strings.HasSuffix(startBlockLines[i], "}") {
					open -= 1
				}
				if open <= 0 {
					break
				}
			}
			if open > 0 {
				err = fmt.Errorf("nginx block '%s' without closed bracket '}'", blockContentLines[blockLineIndex])
				return err
			}
			subBlock, err := parseNginxBlock(blockContentLines[blockLineIndex : blockLineIndex+i+1])
			if err != nil {
				return err
			}
			blockContent.Blocks = append(blockContent.Blocks, subBlock)
			blockLineIndex += i
			continue
		}

		if strings.HasSuffix(line, "}") {
			err = fmt.Errorf("over closed bracket '}' in nginx configuration")
			return err
		}
	}

	return err
}

func parseNginxBlock(blockLines []string) (block BlockT, err error) {
	// Parse name and parameters
	nameParamsParts := strings.Fields(blockLines[0])
	block.Name = nameParamsParts[0]
	block.Params = strings.Join(nameParamsParts[1:len(nameParamsParts)-1], " ")

	// Parse block content
	err = parseNginxBlockContent(&block.BlockContent, blockLines[1:len(blockLines)-1])
	return block, err
}

// Encode functions

func (e *NginxT) EncodeConfigString() (configStr string) {
	configStr = encodeNginxBlockContent(e.ConfigStruct, 0)
	return configStr
}

func encodeNginxBlockContent(blockContent BlockContentT, indent int) (configStr string) {
	indentStr := ""
	for i := 0; i < indent; i++ {
		indentStr += "    "
	}

	for _, val := range blockContent.Directives {
		configStr += indentStr + val.Name + " " + val.Param + " " + val.Value + ";\n"
	}

	for _, val := range blockContent.Blocks {
		configStr += indentStr + val.Name + " " + val.Params + " {\n"
		configStr += encodeNginxBlockContent(val.BlockContent, indent+1)
		configStr += indentStr + "}\n"
	}

	return configStr
}
