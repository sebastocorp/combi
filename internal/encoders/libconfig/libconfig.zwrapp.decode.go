package libconfig

import (
	"fmt"
)

func decodeSettingsTokens(ts []TokenT, cfg map[string]any) (err error) {
	tsLen := len(ts)
	index := 0
	for index < tsLen {
		if ts[index].ttype == TOKEN_TYPE_SEPARATOR {
			index++
			continue
		}

		if ts[index].ttype != TOKEN_TYPE_NAME {
			return fmt.Errorf("malformed configuration, token '%s' not a setting name", ts[index].token)
		}

		valueIndex := index + 2
		if valueIndex > tsLen {
			return fmt.Errorf("malformed configuration, not value in '%s' setting", ts[index].token)
		}

		if ts[index+1].ttype != TOKEN_TYPE_EQUAL {
			return fmt.Errorf("malformed configuration, not equal sign (':' or '=') in '%s' setting", ts[index].token)
		}

		switch ts[valueIndex].ttype {
		case TOKEN_TYPE_VALUE:
			{
				cfg[ts[index].token] = ts[valueIndex].token
				index = valueIndex + 1
				continue
			}
		case TOKEN_TYPE_BRACKET_OPEN:
			{
				offset, err := getScopeOffset(ts[valueIndex:], TOKEN_TYPE_BRACKET_OPEN)
				if err != nil {
					return err
				}

				cfg[ts[index].token], err = decodeArrayTokens(ts[valueIndex+1 : valueIndex+offset])
				if err != nil {
					return err
				}

				index = valueIndex + offset + 1
				continue
			}
		case TOKEN_TYPE_PAREN_OPEN:
			{
				offset, err := getScopeOffset(ts[valueIndex:], TOKEN_TYPE_PAREN_OPEN)
				if err != nil {
					return err
				}

				cfg[ts[index].token], err = decodeListTokens(ts[valueIndex+1 : valueIndex+offset])
				if err != nil {
					return err
				}

				index = valueIndex + offset + 1
				continue
			}
		case TOKEN_TYPE_BRACE_OPEN:
			{
				offset, err := getScopeOffset(ts[valueIndex:], TOKEN_TYPE_BRACE_OPEN)
				if err != nil {
					return err
				}

				cfg[ts[index].token] = make(map[string]any)
				err = decodeSettingsTokens(ts[valueIndex+1:valueIndex+offset], cfg[ts[index].token].(map[string]any))
				if err != nil {
					return err
				}

				index = valueIndex + offset + 1
				continue
			}
		default:
			{
				return fmt.Errorf("malformed configuration, token '%s' unparsable", ts[index+2].token)
			}
		}
	}

	return nil
}

func decodeArrayTokens(ts []TokenT) (array []any, err error) {
	tsLen := len(ts)
	index := 0
	for index < tsLen {
		if ts[index].ttype != TOKEN_TYPE_VALUE {
			return array, fmt.Errorf("malformed configuration, not scalar value '%s' in array", ts[index].token)
		}
		array = append(array, ts[index].token)

		index++
		if index < tsLen {
			if !isSeparator(ts[index].token[0]) {
				return array, fmt.Errorf("malformed configuration, not proper separator '%s' in array", ts[index].token)
			}
			index++
		}
	}

	return array, err
}

func decodeListTokens(ts []TokenT) (list []any, err error) {
	tsLen := len(ts)
	index := 0
	for index < tsLen {
		switch ts[index].ttype {
		case TOKEN_TYPE_VALUE:
			{
				list = append(list, ts[index].token)
			}
		case TOKEN_TYPE_BRACKET_OPEN:
			{
				offset, err := getScopeOffset(ts[index:], TOKEN_TYPE_BRACKET_OPEN)
				if err != nil {
					return list, err
				}

				array, err := decodeArrayTokens(ts[index+1 : index+offset])
				if err != nil {
					return list, err
				}
				list = append(list, array)

				index += offset
			}
		case TOKEN_TYPE_PAREN_OPEN:
			{
				offset, err := getScopeOffset(ts[index:], TOKEN_TYPE_PAREN_OPEN)
				if err != nil {
					return list, err
				}

				subList, err := decodeListTokens(ts[index+1 : index+offset])
				if err != nil {
					return list, err
				}
				list = append(list, subList)

				index += offset
			}
		case TOKEN_TYPE_BRACE_OPEN:
			{
				offset, err := getScopeOffset(ts[index:], TOKEN_TYPE_BRACE_OPEN)
				if err != nil {
					return list, err
				}

				group := make(map[string]any)
				err = decodeSettingsTokens(ts[index+1:index+offset], group)
				if err != nil {
					return list, err
				}
				list = append(list, group)

				index += offset
			}
		default:
			{
				return list, fmt.Errorf("malformed configuration, not proper value '%s' in list, must be a scalar, array, list or group", ts[index].token)
			}
		}

		index++
		if index < tsLen {
			if !isSeparator(ts[index].token[0]) {
				return list, fmt.Errorf("malformed configuration, not proper separator '%s' in list", ts[index].token)
			}
			index++
		}
	}

	return list, err
}

func getScopeOffset(ts []TokenT, openScopeType TokenTypeT) (offset int, err error) {
	tsLen := len(ts)
	open := 0
	for ; offset < tsLen; offset++ {
		if ts[offset].ttype == openScopeType {
			open++
		}
		if ts[offset].ttype == openScopeType+1 {
			open--
			if open == 0 {
				break
			}
		}
	}

	if open != 0 {
		err = fmt.Errorf("malformed configuration, scope not closed")
	}

	return offset, err
}
