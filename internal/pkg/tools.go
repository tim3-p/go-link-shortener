package pkg

import (
	"crypto/md5"
	"encoding/hex"
)

func HashURL(url []byte) string {
	hash := md5.Sum([]byte(url))
	str := hex.EncodeToString(hash[:])
	return string(str[1:6])
}