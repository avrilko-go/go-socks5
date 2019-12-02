package util

import (
	"encoding/base64"
)

var password = "fRB+remQqTY/OlXNa6THL5Ij4yK3QlJM/t3uhv9dxgZGgvDa5mHPVgJf3nRtOE8eMyWah7F7QKG8zH81VBuo2HygGPF49r/DKdKvi3IMMAvRjW+0hWy4mbka7PQIZfzo1ABaXnX6R2STs9xmKEWj1z08wFhcH2eBmxJ5S8p64KsR1oynIZHly56yIIjtsFNQdrZNOwoWQTf9n3DVnarZFFssLURz8lGmF6Uc+3EVSUMOgOvh94ROrAVI37u6l/kZxQ0EbjGP0JYPJsIr28iJjjljSlk04r7BtYouosl3Ewcy7ycDKr0JJJVXlB1orvjzac7nYsQ+6phg5PWcAYNq0w=="

var passwordOrigin = make([]byte, 256)

var passwordEncode = make([]byte, 256)

func init() {
	originByte, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		panic(originByte)
	}
	passwordOrigin = originByte
	encodeByte := make([]byte, 256)
	// 生成加密串
	for k, v := range originByte {
		encodeByte[v] = byte(k)
	}
	passwordEncode = encodeByte
}

func Encode(b []byte) {
	for k, v := range b {
		b[k] = passwordEncode[v]
	}
}

func Decode(b []byte) {
	for k, v := range b {
		b[k] = passwordOrigin[v]
	}
}
