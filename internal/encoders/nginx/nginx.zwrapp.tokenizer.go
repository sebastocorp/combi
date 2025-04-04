package nginx

import (
	"fmt"
	"slices"
)

const (
	TOKEN_TYPE_DEFAULT TokenTypeT = iota
	TOKEN_TYPE_ITEM
	TOKEN_TYPE_SEPARATOR   // ';'
	TOKEN_TYPE_BLOCK_OPEN  // '{'
	TOKEN_TYPE_BLOCK_CLOSE // '}'
)

type TokenTypeT int

type TokenT struct {
	ttype TokenTypeT
	value string
}

func tokenize(cfgBytes []byte) (tokens []TokenT, err error) {
	cfgBytesLen := len(cfgBytes)
	index := 0
	for index < cfgBytesLen {
		if isSpace(cfgBytes[index]) {
			index++
			continue
		}

		if cfgBytes[index] == '#' {
			index += getOffsetUntil(cfgBytes[index:], '\n')
			continue
		}

		if isSeparator(cfgBytes[index]) {
			tokens = append(tokens, TokenT{
				ttype: TOKEN_TYPE_SEPARATOR,
				value: string(cfgBytes[index : index+1]),
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
		index += len(t.value)
	}

	return tokens, err
}

func isSpace(b byte) bool {
	return slices.Contains([]byte{' ', '\t', '\n'}, b)
}

func isSeparator(b byte) bool {
	return slices.Contains([]byte{';'}, b)
}

func isScope(b byte) bool {
	return slices.Contains([]byte{'{', '}'}, b)
}

func isSpecialToken(b byte) bool {
	return slices.Contains([]byte{' ', '\t', '\n', ';', '{', '}', '#'}, b)
}

// func getTokenTypeString(tt TokenTypeT) string {
// 	return [...]string{
// 		"TOKEN_TYPE_DEFAULT",
// 		"TOKEN_TYPE_ITEM",
// 		"TOKEN_TYPE_SEPARATOR",
// 		"TOKEN_TYPE_BLOCK_OPEN",
// 		"TOKEN_TYPE_BLOCK_CLOSE",
// 	}[tt]
// }

func getScopeToken(b byte) (t TokenT) {
	t.value = string([]byte{b})
	switch b {
	case '{':
		t.ttype = TOKEN_TYPE_BLOCK_OPEN
	case '}':
		t.ttype = TOKEN_TYPE_BLOCK_CLOSE
	}

	return t
}

func getToken(cfgBytes []byte) (t TokenT, err error) {
	cfgBytesLen := len(cfgBytes)
	t.ttype = TOKEN_TYPE_ITEM
	if cfgBytes[0] == '"' || cfgBytes[0] == '\'' {
		strDelim := cfgBytes[0]
		i := 1
		for ; i < cfgBytesLen; i++ {
			if cfgBytes[i] == strDelim && cfgBytes[i] != '\\' {
				i++
				break
			}
		}

		t.value = string(cfgBytes[:i])
		if t.value[len(t.value)-1] != strDelim {
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

	t.value = string(cfgBytes[:i])

	return t, err
}

func getOffsetUntil(cfgBytes []byte, b byte) (offset int) {
	tsLen := len(cfgBytes)
	offset = 0
	for ; offset < tsLen; offset++ {
		if cfgBytes[offset] == b {
			break
		}
	}

	return offset
}
