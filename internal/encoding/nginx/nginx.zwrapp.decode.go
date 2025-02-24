package nginx

import "fmt"

func decodeConfigTokens(ts []TokenT, cfg map[string]any) (err error) {
	tsLen := len(ts)
	for ti := 0; ti < tsLen; ti++ {
		if ts[ti].ttype == TOKEN_TYPE_ITEM {
			cfg[ts[ti].value] = nil

			var value string
			var tt TokenTypeT
			var offset int
			value, tt, offset, err = getOffsetUntilNextSpecialToken(ts[ti:])
			if err != nil {
				return err
			}
			fmt.Printf("--------------------------------------------------\n")
			fmt.Printf("type: %-24s; len: %-5d; item: %-10s\n", getTokenTypeString(ts[ti].ttype), len(ts[ti].value), ts[ti].value)
			fmt.Printf("type: %-24s; len: %-5d; item: %s\n", getTokenTypeString(tt), len(value), value)
			fmt.Printf("type: %-24s; len: %-5d; item: %-10s\n", getTokenTypeString(ts[ti+offset].ttype), len(ts[ti+offset].value), ts[ti+offset].value)
			fmt.Printf("--------------------------------------------------\n")
			ti += offset
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

	if tt == TOKEN_TYPE_DEFAULT {
		err = fmt.Errorf("malformed configuration, unable to now next special token")
	}

	if len(value) > 0 {
		value = value[:len(value)-1]
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
