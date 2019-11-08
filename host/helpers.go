package host

import (
	"errors"
	"github.com/libp2p/go-libp2p-core/crypto"
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

// RoutingStateFromHost returns a routing.RoutingState record that contains the Host's
// ID and all of its listen Addrs.
func RoutingStateFromHost(h Host) *routing.RoutingState {
	return routing.RoutingStateFromAddrInfo(InfoFromHost(h))
}

// SignedRoutingStateFromHost
func SignedRoutingStateFromHost(h Host) (*crypto.SignedEnvelope, error) {
	privKey := h.Peerstore().PrivKey(h.ID())
	if privKey == nil {
		return nil, errors.New("unable to find host's private key in peerstore")
	}

	return RoutingStateFromHost(h).ToSignedEnvelope(privKey)
}
