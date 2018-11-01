package securitys

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
)

func GetHmacCode(source, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	io.WriteString(h, source)
	return fmt.Sprintf("%x", h.Sum(nil))
}
