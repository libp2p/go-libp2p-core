package event

import (
	peer "github.com/libp2p/go-libp2p-core/peer"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
)

// EvtPeerProtocolsUpdated should be emitted when a peer we're connected to adds or removes protocols from their stack.
type EvtPeerProtocolsUpdated struct {
	// Peer is the peer whose protocols were updated.
	Peer peer.ID
	// Added enumerates the protocols that were added by this peer.
	Added []protocol.ID
	// Removed enumerates the protocols that were removed by this peer.
	Removed []protocol.ID
}
