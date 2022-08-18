// Deprecated: This package has moved into go-libp2p as a sub-package: github.com/libp2p/go-libp2p/core/sec/insecure.
//
// Package insecure provides an insecure, unencrypted implementation of the the SecureConn and SecureTransport interfaces.
//
// Recommended only for testing and other non-production usage.
package insecure

import (
	ci "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/sec/insecure"
)

// ID is the multistream-select protocol ID that should be used when identifying
// this security transport.
// Deprecated: use github.com/libp2p/go-libp2p/core/sec/insecure.ID instead
const ID = insecure.ID

// Transport is a no-op stream security transport. It provides no
// security and simply mocks the security methods. Identity methods
// return the local peer's ID and private key, and whatever the remote
// peer presents as their ID and public key.
// No authentication of the remote identity is performed.
// Deprecated: use github.com/libp2p/go-libp2p/core/sec/insecure.Transport instead
type Transport = insecure.Transport

// NewWithIdentity constructs a new insecure transport. The provided private key
// is stored and returned from LocalPrivateKey to satisfy the
// SecureTransport interface, and the public key is sent to
// remote peers. No security is provided.
// Deprecated: use github.com/libp2p/go-libp2p/core/sec/insecure.NewWithIdentity instead
func NewWithIdentity(id peer.ID, key ci.PrivKey) *Transport {
	return insecure.NewWithIdentity(id, key)
}

// Conn is the connection type returned by the insecure transport.
// Deprecated: use github.com/libp2p/go-libp2p/core/sec/insecure.Conn instead
type Conn = insecure.Conn
