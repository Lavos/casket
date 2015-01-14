package casket

import (
	"fmt"
	"crypto/sha1"
	"encoding/hex"
)

type SHA1Sum [20]byte

func NewSHA1Sum (content []byte) SHA1Sum {
	return SHA1Sum(sha1.Sum(content))
}

func NewSHA1SumFromString (s string) SHA1Sum {
	var n SHA1Sum
	b, _ := hex.DecodeString(s)

	copy(n[:], b)

	return n
}

func NewSHA1SumFromBytes (b []byte) SHA1Sum {
	var n SHA1Sum
	copy(n[:], b)
	return n
}

func (s SHA1Sum) String() string {
	return hex.EncodeToString(s[:])
}

func (s SHA1Sum) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s.String())), nil
}
