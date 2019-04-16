package host

import "github.com/libp2p/go-libp2p-core/peer"

// InfoFromHost returns a peer.Info struct with the Host's ID and all of its Addrs.
func InfoFromHost(h Host) *peer.Info {
	return &peer.Info{
		ID:    h.ID(),
		Addrs: h.Addrs(),
	}
}
