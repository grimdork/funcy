package funcy

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// GenerateToken creates a token for the user session.
func GenerateToken(username string) string {
	return Base6424(GenerateString(64) + username)
}

const validChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!'#$%&/()=?@*^<>-.:,;|[]{}"

// GenerateString of length n.
func GenerateString(n int) string {
	s := make([]byte, n)
	for i := 0; i < n; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(validChars))))
		if err != nil {
			return ""
		}
		c := validChars[n.Int64()]
		s[i] = c
	}
	return string(s)
}

const alphabet = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Base6424 is used by some password hashing algorithms.
func Base6424(src string) string {
	if len(src) == 0 {
		return ""
	}

	buf := strings.Builder{}
	for len(src) > 0 {
		switch len(src) {
		default:
			_ = buf.WriteByte(alphabet[src[0]&0x3f])
			_ = buf.WriteByte(alphabet[((src[0]>>6)|(src[1]<<2))&0x3f])
			_ = buf.WriteByte(alphabet[((src[1]>>4)|(src[2]<<4))&0x3f])
			_ = buf.WriteByte(alphabet[(src[2]>>2)&0x3f])
			src = src[3:]
		case 2:
			_ = buf.WriteByte(alphabet[src[0]&0x3f])
			_ = buf.WriteByte(alphabet[((src[0]>>6)|(src[1]<<2))&0x3f])
			_ = buf.WriteByte(alphabet[(src[1]>>4)&0x3f])
			src = src[2:]
		case 1:
			_ = buf.WriteByte(alphabet[src[0]&0x3f])
			_ = buf.WriteByte(alphabet[(src[0]>>6)&0x3f])
			src = src[1:]
		}
	}
	return buf.String()
}
