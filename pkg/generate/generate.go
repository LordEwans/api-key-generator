package generate

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func GenerateKey() string {
	seed := make([]byte, 16)
	io.ReadFull(rand.Reader, seed)

	return hex.EncodeToString(seed)
}
