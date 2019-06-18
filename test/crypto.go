package test

import (
	"math/rand"
	"sync/atomic"
	"time"

	ci "github.com/libp2p/go-libp2p-core/crypto"
)

var generatedPairs int64 = 0

func RandTestKeyPair(typ, bits int) (ci.PrivKey, ci.PubKey, error) {
	seed := time.Now().UnixNano()

	// workaround for low time resolution
	seed += atomic.AddInt64(&generatedPairs, 1) << 32

	return SeededTestKeyPair(typ, bits, time.Now().UnixNano())
}

func SeededTestKeyPair(typ, bits int, seed int64) (ci.PrivKey, ci.PubKey, error) {
	r := rand.New(rand.NewSource(seed))
	return ci.GenerateKeyPairWithReader(typ, bits, r)
}
