package identify

import (
	peer "github.com/libp2p/go-libp2p-core/peer"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
)

type PeerProtocolsUpdated struct {
	Peer    peer.ID
	Added   []protocol.ID
	Removed []protocol.ID
}
