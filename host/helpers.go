package host

import (
	"errors"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/routing"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
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
//
// By default, only publicly routable addresses will be included.
// To include loopback and LAN addresses, pass in the IncludeLocalAddrs option:
//
//    state := SignedRoutingStateFromHost(h, IncludeLocalAddrs)
func SignedRoutingStateFromHost(h minimalHost, opts ...Option) (*routing.SignedRoutingState, error) {
	cfg := config{}
	for _, opt := range opts {
		opt(&cfg)
	}

	privKey := h.Peerstore().PrivKey(h.ID())
	if privKey == nil {
		return nil, errors.New("unable to find host's private key in peerstore")
	}

	var addrs []ma.Multiaddr
	if cfg.includeLocalAddrs {
		addrs = h.Addrs()
	} else {
		for _, a := range h.Addrs() {
			if manet.IsPublicAddr(a) {
				addrs = append(addrs, a)
			}
		}
	}

	return routing.MakeSignedRoutingState(privKey, addrs)
}

// IncludeLocalAddrs can be passed into SignedRoutingStateFromHost to
// produce a routing record with LAN and loopback addresses included.
func IncludeLocalAddrs(cfg *config) {
	cfg.includeLocalAddrs = true
}

// minimalHost is the subset of the Host interface that's required by
// SignedRoutingStateFromHost.
type minimalHost interface {
	ID() peer.ID
	Peerstore() peerstore.Peerstore
	Addrs() []ma.Multiaddr
}

type Option func(cfg *config)
type config struct {
	includeLocalAddrs bool
}
