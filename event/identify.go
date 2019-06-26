package event

import "github.com/libp2p/go-libp2p-core/peer"

// EvtPeerInitialIdentificationCompleted is emitted when the initial identification round for a peer is completed.
type EvtPeerInitialIdentificationCompleted struct {
	// Peer is the ID of the peer whose identification succeeded.
	Peer peer.ID
}

// EvtPeerInitialIdentificationFailed is emitted when the initial identification round for a peer failed.
type EvtPeerInitialIdentificationFailed struct {
	// Peer is the ID of the peer whose identification failed.
	Peer peer.ID
	// Reason is the reason why identification failed.
	Reason error
}
