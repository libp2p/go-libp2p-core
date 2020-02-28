package event

import "github.com/libp2p/go-libp2p-core/network"

// EvtPeerConnectionStateChange should be emitted when we connect to a peer or disconnect
// from a peer. It contains the network interface for the connection,
// the connection handle & the new state of the connection.
type EvtPeerConnectionStateChange struct {
	Network    network.Network
	Connection network.Conn
	NewState   network.Connectedness
}

// EvtStreamStateChange should be emitted when we open a new stream with a peer or close an existing stream.
// It contains the network interface for the connection, the stream handle &
// the new state of the stream.
type EvtStreamStateChange struct {
	Network  network.Network
	Stream   network.Stream
	NewState network.Connectedness
}
