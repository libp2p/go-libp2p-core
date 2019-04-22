package test

import (
	"math/rand"
	"time"

	ci "github.com/libp2p/go-libp2p-core/crypto"
)

func RandTestKeyPair(typ, bits int) (ci.PrivKey, ci.PubKey, error) {
	return SeededTestKeyPair(typ, bits, time.Now().UnixNano())
}

func SeededTestKeyPair(typ, bits int, seed int64) (ci.PrivKey, ci.PubKey, error) {
	r := rand.New(rand.NewSource(seed))
	return ci.GenerateKeyPairWithReader(typ, bits, r)
}
