package mfa

import (
	"encoding/base32"
	"strings"
)

var (
	base32LookupTable = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567")
)

func base32Decode(secret string) ([]byte, error) {
	secret = strings.ToUpper(strings.TrimRight(secret, "="))
	return base32.StdEncoding.DecodeString(secret)
}
