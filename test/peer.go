package test

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/test"
)

// Deprecated: use github.com/libp2p/go-libp2p/core/test.RandPeerID instead
func RandPeerID() (peer.ID, error) {
	return test.RandPeerID()
}

// Deprecated: use github.com/libp2p/go-libp2p/core/test.RandPeerIDFatal instead
func RandPeerIDFatal(t testing.TB) peer.ID {
	return test.RandPeerIDFatal(t)
}
