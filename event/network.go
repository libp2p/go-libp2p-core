package event

import (
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

// EvtPeerConnectednessChanged should be emitted every time we form a connection with a peer or drop our last
// connection with the peer. Essentially, it is emitted in two cases:
// a) We form a/any connection with a peer.
// b) We go from having a connection/s with a peer to having no connection with the peer.
// It contains the Id of the remote peer and the new connectedness state.
type EvtPeerConnectednessChanged struct {
	Peer          peer.ID
	Connectedness network.Connectedness
}
