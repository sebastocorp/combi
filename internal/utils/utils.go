package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func GenHashString(args ...string) (h string, err error) {
	str := ""
	for _, sv := range args {
		str += sv
	}

	md5Hash := md5.New()
	_, err = md5Hash.Write([]byte(str))
	if err != nil {
		return h, err
	}

	h = hex.EncodeToString(md5Hash.Sum(nil))

	return h, err
}
