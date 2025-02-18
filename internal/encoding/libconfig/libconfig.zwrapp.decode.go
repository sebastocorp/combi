package libconfig

import (
	"fmt"
	"regexp"
	"slices"
)

const (
	TOKEN_TYPE_NAME TokenTypeT = iota
	TOKEN_TYPE_VALUE
	TOKEN_TYPE_EQUAL         // ':', '='
	TOKEN_TYPE_SEPARATOR     // ';', ','
	TOKEN_TYPE_BRACKET_OPEN  // '['
	TOKEN_TYPE_BRACKET_CLOSE // ']'
	TOKEN_TYPE_PAREN_OPEN    // ')'
	TOKEN_TYPE_PAREN_CLOSE   // '('
	TOKEN_TYPE_BRACE_OPEN    // '{'
	TOKEN_TYPE_BRACE_CLOSE   // '}'
)

var (
	tokenTypeRegexName    = regexp.MustCompile(`^[A-Za-z][-A-Za-z0-9_]*$`)
	tokenTypeRegexBool    = regexp.MustCompile(`^([Tt][Rr][Uu][Ee])|([Ff][Aa][Ll][Ss][Ee])$`)
	tokenTypeRegexInteger = regexp.MustCompile(`^[-+]?[0-9]+(L(L)?)?$`)
	tokenTypeRegexHex     = regexp.MustCompile(`^0[Xx][0-9A-Fa-f]+(L(L)?)?$`)
	tokenTypeRegexFloat   = regexp.MustCompile(`^([-+]?([0-9]*)?.[0-9]*([eE][-+]?[0-9]+)?)|([-+]([0-9]+)(.[0-9]*)?[eE][-+]?[0-9]+)$`)
)

type TokenTypeT int

type TokenT struct {
	ttype TokenTypeT
	token string
}

func tokenize(cfgBytes []byte) (tokens []TokenT, err error) {
	cfgBytesLen := len(cfgBytes)
	index := 0
	for index < cfgBytesLen {
		if isSpace(cfgBytes[index]) {
			index++
			continue
		}

		if isEqual(cfgBytes[index]) {
			tokens = append(tokens, TokenT{
				ttype: TOKEN_TYPE_EQUAL,
				token: string(cfgBytes[index : index+1]),
			})
			index++
			continue
		}

		if isSeparator(cfgBytes[index]) {
			tokens = append(tokens, TokenT{
				ttype: TOKEN_TYPE_SEPARATOR,
				token: string(cfgBytes[index : index+1]),
			})
			index++
			continue
		}

		if isScope(cfgBytes[index]) {
			tokens = append(tokens, getScopeToken(cfgBytes[index]))
			index++
			continue
		}

		t, err := getToken(cfgBytes[index:])
		if err != nil {
			return tokens, err
		}
		tokens = append(tokens, t)
		index += len(t.token)
	}

	return tokens, err
}

func isSpace(b byte) bool {
	return slices.Contains([]byte{' ', '\t', '\n'}, b)
}

func isEqual(b byte) bool {
	return slices.Contains([]byte{':', '='}, b)
}

func isSeparator(b byte) bool {
	return slices.Contains([]byte{';', ','}, b)
}

func isScope(b byte) bool {
	return slices.Contains([]byte{'[', ']', '(', ')', '{', '}'}, b)
}

func getScopeToken(b byte) (t TokenT) {
	t.token = string([]byte{b})
	switch b {
	case '[':
		t.ttype = TOKEN_TYPE_BRACKET_OPEN
	case ']':
		t.ttype = TOKEN_TYPE_BRACKET_CLOSE
	case '{':
		t.ttype = TOKEN_TYPE_BRACE_OPEN
	case '}':
		t.ttype = TOKEN_TYPE_BRACE_CLOSE
	case '(':
		t.ttype = TOKEN_TYPE_PAREN_OPEN
	case ')':
		t.ttype = TOKEN_TYPE_PAREN_CLOSE
	}

	return t
}

func getToken(cfgBytes []byte) (t TokenT, err error) {
	cfgBytesLen := len(cfgBytes)

	if cfgBytes[0] == '"' {
		t.ttype = TOKEN_TYPE_VALUE
		i := 1
		for ; i < cfgBytesLen; i++ {
			if cfgBytes[i] == '"' && cfgBytes[i] != '\\' {
				i++
				break
			}
		}

		t.token = string(cfgBytes[:i])

		if t.token[len(t.token)-1] != '"' {
			err = fmt.Errorf("unclosed string")
		}

		return t, err
	}

	i := 1
	for ; i < cfgBytesLen; i++ {
		if isSpace(cfgBytes[i]) || isEqual(cfgBytes[i]) || isSeparator(cfgBytes[i]) || isScope(cfgBytes[i]) {
			break
		}
	}

	t.token = string(cfgBytes[:i])

	t.ttype, err = getTypeFromToken(t)

	return t, err
}

func getTypeFromToken(t TokenT) (tt TokenTypeT, err error) {
	tt = TOKEN_TYPE_NAME

	if !tokenTypeRegexName.MatchString(t.token) {
		if tokenTypeRegexBool.MatchString(t.token) ||
			tokenTypeRegexInteger.MatchString(t.token) ||
			tokenTypeRegexHex.MatchString(t.token) ||
			tokenTypeRegexFloat.MatchString(t.token) {
			tt = TOKEN_TYPE_VALUE
		} else {
			err = fmt.Errorf("unrecognize '%s' token type", t.token)
		}
	}

	return tt, err
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

func encodeSettingsTokens(ts []TokenT, cfg map[string]any) (err error) {
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

				cfg[ts[index].token], err = encodeArrayTokens(ts[valueIndex+1 : valueIndex+offset])
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

				cfg[ts[index].token], err = encodeListTokens(ts[valueIndex+1 : valueIndex+offset])
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
				err = encodeSettingsTokens(ts[valueIndex+1:valueIndex+offset], cfg[ts[index].token].(map[string]any))
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

func encodeArrayTokens(ts []TokenT) (array []any, err error) {
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

func encodeListTokens(ts []TokenT) (list []any, err error) {
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

				array, err := encodeArrayTokens(ts[index+1 : index+offset])
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

				subList, err := encodeListTokens(ts[index+1 : index+offset])
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
				err = encodeSettingsTokens(ts[index+1:index+offset], group)
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
