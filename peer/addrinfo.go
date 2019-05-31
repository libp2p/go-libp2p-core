package peer

import (
	"fmt"

	ma "github.com/multiformats/go-multiaddr"
)

// AddrInfo is a small struct used to pass around a peer with
// a set of addresses (and later, keys?).
type AddrInfo struct {
	ID    ID
	Addrs []ma.Multiaddr
}

var _ fmt.Stringer = AddrInfo{}

func (pi AddrInfo) String() string {
	return fmt.Sprintf("{%v: %v}", pi.ID, pi.Addrs)
}

var ErrInvalidAddr = fmt.Errorf("invalid p2p multiaddr")

func AddrInfoFromP2pAddr(m ma.Multiaddr) (*AddrInfo, error) {
	if m == nil {
		return nil, ErrInvalidAddr
	}

	transport, p2ppart := ma.SplitLast(m)
	if p2ppart == nil || p2ppart.Protocol().Code != ma.P_P2P {
		return nil, ErrInvalidAddr
	}
	id, err := IDFromBytes(p2ppart.RawValue())
	if err != nil {
		return nil, err
	}
	info := &AddrInfo{ID: id}
	if transport != nil {
		info.Addrs = []ma.Multiaddr{transport}
	}
	return info, nil
}

func AddrInfoToP2pAddrs(pi *AddrInfo) ([]ma.Multiaddr, error) {
	var addrs []ma.Multiaddr
	p2ppart, err := ma.NewComponent("p2p", IDB58Encode(pi.ID))
	if err != nil {
		return nil, err
	}
	for _, addr := range pi.Addrs {
		addrs = append(addrs, addr.Encapsulate(p2ppart))
	}
	return addrs, nil
}

func (pi *AddrInfo) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"peerID": pi.ID.Pretty(),
		"addrs":  pi.Addrs,
	}
}
