package event

import (
	"github.com/libp2p/go-libp2p-core/record"
	ma "github.com/multiformats/go-multiaddr"
)

// EvtLocalAddressesUpdated should be emitted when the set of listen addresses for
// the local host changes. This may happen for a number of reasons. For example,
// we may have opened a new relay connection, established a new NAT mapping via
// UPnP, or been informed of our observed address by another peer.
type EvtLocalAddressesUpdated struct {
	// Added enumerates the listen addresses that were added for the local peer.
	Added []ma.Multiaddr

	// Removed enumerates listen addresses that were removed from the local peer.
	Removed []ma.Multiaddr
}

// EvtLocalPeerRoutingStateUpdated should be emitted when a new signed PeerRecord
// for the local peer has been produced. This will happen whenever the set of listen
// addresses changes.
type EvtLocalPeerRecordUpdated struct {
	SignedRecord *record.SignedEnvelope
}
