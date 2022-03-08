package util

import (
	"crypto/sha1"
	"encoding/hex"
	"sort"
	"strings"
)

func Verify(signature string, timestamp string, nonce string, token string) bool {
	arr := []string{timestamp, nonce, token}
	sort.Strings(arr)

	n := len(timestamp) + len(nonce) + len(token)
	builder := &strings.Builder{}
	builder.Grow(n)

	for _, s := range arr {
		builder.WriteString(s)
	}

	return encrypt(builder.String()) == signature
}

func encrypt(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
