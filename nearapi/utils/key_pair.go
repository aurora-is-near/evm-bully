package utils

import "crypto/ed25519"

// All supported key types
const (
	ED25519 = 0
)

// PublicKey encoding for NEAR.
type PublicKey struct {
	KeyType uint8
	Data    [32]byte
}

// PublicKeyFromEd25519 derives a public key in NEAR encoding from pk.
func PublicKeyFromEd25519(pk ed25519.PublicKey) PublicKey {
	var pubKey PublicKey
	pubKey.KeyType = ED25519
	copy(pubKey.Data[:], pk)
	return pubKey
}
