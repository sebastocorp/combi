package nginx

import "fmt"

const (
	nginxBlockPatternKey = "nginxBlockPattern"
)

func decodeConfigTokens(ts []TokenT, cfg map[string]any) (err error) {
	tsLen := len(ts)
	for ti := 0; ti < tsLen; ti++ {
		if ts[ti].ttype == TOKEN_TYPE_ITEM {
			var value string
			var tt TokenTypeT
			var offset int
			value, tt, offset, err = getOffsetUntilNextSpecialToken(ts[ti:])
			if err != nil {
				return err
			}

			var scopeOffset int
			if tt == TOKEN_TYPE_BLOCK_OPEN {
				scopeOffset, err = getScopeOffset(ts[ti+offset:], tt)
				if err != nil {
					return err
				}

				if _, ok := cfg[ts[ti].value]; ok {
					new := make(map[string]any)
					new[nginxBlockPatternKey] = value
					err = decodeConfigTokens(ts[ti+offset+1:ti+offset+scopeOffset], new)
					if err != nil {
						return err
					}

					switch cfg[ts[ti].value].(type) {
					case []any:
						{
							cfg[ts[ti].value] = append(cfg[ts[ti].value].([]any), new)
						}
					case map[string]any:
						{
							old := cfg[ts[ti].value]
							cfg[ts[ti].value] = []any{old, new}
						}
					}
				} else {
					cfg[ts[ti].value] = make(map[string]any)
					cfg[ts[ti].value].(map[string]any)[nginxBlockPatternKey] = value
					err = decodeConfigTokens(ts[ti+offset+1:ti+offset+scopeOffset], cfg[ts[ti].value].(map[string]any))
					if err != nil {
						return err
					}
				}
			}

			if tt == TOKEN_TYPE_SEPARATOR {
				if _, ok := cfg[ts[ti].value]; ok {
					switch cfg[ts[ti].value].(type) {
					case []any:
						{
							cfg[ts[ti].value] = append(cfg[ts[ti].value].([]any), value)
						}
					case string:
						{
							old := cfg[ts[ti].value]
							cfg[ts[ti].value] = []any{old, value}
						}
					}
				} else {
					cfg[ts[ti].value] = value
				}
			}

			ti += offset + scopeOffset
		}
	}

	return nil
}

func getOffsetUntilNextSpecialToken(ts []TokenT) (value string, tt TokenTypeT, offset int, err error) {
	tsLen := len(ts)
	tt = TOKEN_TYPE_DEFAULT
	offset = 1
	for ; offset < tsLen; offset++ {
		if ts[offset].ttype == TOKEN_TYPE_ITEM {
			value += ts[offset].value + " "
		}
		if ts[offset].ttype == TOKEN_TYPE_SEPARATOR || ts[offset].ttype == TOKEN_TYPE_BLOCK_OPEN {
			tt = ts[offset].ttype
			break
		}
	}

	if len(value) > 0 {
		value = value[:len(value)-1]
	}

	if tt == TOKEN_TYPE_DEFAULT {
		err = fmt.Errorf("malformed configuration, unable to now next special token")
	}

	return value, tt, offset, err
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
