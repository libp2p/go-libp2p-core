package event

import (
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

// EvtHolePunchConnSuccessful is emitted when a hole punching attempt to establish a direct
// connection with a peer is successful.
type EvtHolePunchConnSuccessful struct {
	// Peer is the ID of the peer we hole punched with
	Peer peer.ID
	// ProxyConn is the proxy connection over which we co-ordinated hole punching.
	// In the current implementation, this connection will be closed after a grace period.
	// It is the user's responsibility to migrate all streams from this connection to the new hole
	// punched connection.
	ProxyConn network.Conn
}
