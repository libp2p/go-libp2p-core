package test

import (
	ci "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/test"
)

// Deprecated: use github.com/libp2p/go-libp2p/core/sec.RandTestKeyPair instead
func RandTestKeyPair(typ, bits int) (ci.PrivKey, ci.PubKey, error) {
	return test.RandTestKeyPair(typ, bits)
}

// Deprecated: use github.com/libp2p/go-libp2p/core/sec.SeededTestKeyPair instead
func SeededTestKeyPair(typ, bits int, seed int64) (ci.PrivKey, ci.PubKey, error) {
	return test.SeededTestKeyPair(typ, bits, seed)
}
