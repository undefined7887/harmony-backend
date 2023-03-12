package domain

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/undefined7887/harmony-backend/internal/util"
	cryptoutil "github.com/undefined7887/harmony-backend/internal/util/crypto"
)

const (
	IdSize = sha256.Size224

	TokenSize = 64
)

var Base64 = base64.RawURLEncoding

func ID() string {
	return Base64.EncodeToString(cryptoutil.Bytes(IdSize))
}

func CombineIDs(ids ...string) string {
	// Soring ids lexicographically
	util.Sort(ids)

	// Combining all ids
	result := make([]byte, len(ids)*IdSize)

	for i, id := range ids {
		if _, err := Base64.Decode(result[IdSize*i:IdSize*(i+1)], []byte(id)); err != nil {
			panic(fmt.Sprintf("error combining ids: decode base64: %v", err))
		}
	}

	return base64.RawURLEncoding.EncodeToString(cryptoutil.Sha224(result))
}

func Token() string {
	return base64.RawURLEncoding.EncodeToString(cryptoutil.Bytes(TokenSize))
}
