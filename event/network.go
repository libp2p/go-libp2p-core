package event

import "github.com/libp2p/go-libp2p-core/network"

// EvtPeerStateChange should be emitted when we form our first connection with a peer or drop our last
// connection with the peer. Essentially, it is emitted in two cases:
// a) We go from having no connection with a peer to having a connection with a peer.
// b) We go from having a connection/s with a peer to having no connection with the peer.
// It contains the network interface for the connection, the connection handle for the first/last connection and
// the new connection state.
type EvtPeerStateChange struct {
	Connection network.Conn
	NewState   network.Connectedness
}
