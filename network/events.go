package network

// EvtPeerConnectionStateChange should be emitted when we connect/disconnect from a peer
type EvtPeerConnectionStateChange struct {
	Network    Network
	Connection Conn
	NewState   Connectedness
}

// EvtStreamStateChange is emitted when we open/close a stream with a peer
type EvtStreamStateChange struct {
	Network  Network
	Stream   Stream
	NewState Connectedness
}
