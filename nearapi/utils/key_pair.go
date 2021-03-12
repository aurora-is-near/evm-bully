package utils

import "crypto/ed25519"

// All supported key types
const (
	ED25519 = 0
)

type PublicKey struct {
	KeyType int
	Data    []byte
}

func PublicKeyFromEd25519(pk ed25519.PublicKey) *PublicKey {
	var pubKey PublicKey
	pubKey.KeyType = ED25519
	pubKey.Data = pk
	return &pubKey
}
