package libconfig

import (
	"fmt"
	"slices"
)

const (
	TOKEN_TYPE_UNDEFINED TokenTypeT = iota
	TOKEN_TYPE_NAME
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

type TokenTypeT int

type TokenT struct {
	ttype TokenTypeT
	token string
}

func tokenize(cfgBytes []byte) (ts []TokenT, err error) {
	cfgBytesLen := len(cfgBytes)
	index := 0
	for index < cfgBytesLen {
		if isSpace(cfgBytes[index]) {
			index++
			continue
		}

		if isEqual(cfgBytes[index]) {
			ts = append(ts, TokenT{
				ttype: TOKEN_TYPE_EQUAL,
				token: string(cfgBytes[index : index+1]),
			})
			index++
			continue
		}

		if isSeparator(cfgBytes[index]) {
			ts = append(ts, TokenT{
				ttype: TOKEN_TYPE_SEPARATOR,
				token: string(cfgBytes[index : index+1]),
			})
			index++
			continue
		}

		if isScope(cfgBytes[index]) {
			ts = append(ts, getScopeToken(cfgBytes[index]))
			index++
			continue
		}

		t, err := getToken(cfgBytes[index:])
		if err != nil {
			return ts, err
		}

		if t.ttype == TOKEN_TYPE_UNDEFINED {
			tsIndex := len(ts) - 1

			t.ttype = TOKEN_TYPE_NAME
			if tsIndex >= 0 && ts[tsIndex].ttype == TOKEN_TYPE_EQUAL {
				t.ttype = TOKEN_TYPE_VALUE
			}
		}
		ts = append(ts, t)

		index += len(t.token)
	}

	return ts, err
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
		if isSpecialToken(cfgBytes[i]) {
			break
		}
	}

	t.token = string(cfgBytes[:i])

	return t, err
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

func isSpecialToken(b byte) bool {
	return slices.Contains([]byte{'[', ']', '(', ')', '{', '}', ';', ',', ':', '=', ' ', '\t', '\n'}, b)
}
