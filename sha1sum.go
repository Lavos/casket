package casket

import (
	"crypto/sha1"
	"encoding/hex"
)

type SHA1Sum [20]byte

func NewSHA1Sum (content []byte) SHA1Sum {
	return SHA1Sum(sha1.Sum(content))
}

func (s SHA1Sum) String() string {
	return hex.EncodeToString(s[:])
}
