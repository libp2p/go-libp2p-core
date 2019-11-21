package host

import (
	"errors"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
)

// InfoFromHost returns a peer.AddrInfo struct with the Host's ID and all of its Addrs.
func InfoFromHost(h Host) *peer.AddrInfo {
	return &peer.AddrInfo{
		ID:    h.ID(),
		Addrs: h.Addrs(),
	}
}

// SignedRoutingStateFromHost returns a SignedRoutingState record containing
// the Host's listen addresses, signed with the Host's private key.
func SignedRoutingStateFromHost(h Host) (*routing.SignedRoutingState, error) {
	privKey := h.Peerstore().PrivKey(h.ID())
	if privKey == nil {
		return nil, errors.New("unable to find host's private key in peerstore")
	}

	return routing.MakeSignedRoutingState(privKey, h.Addrs())
}
