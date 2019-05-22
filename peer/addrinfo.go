package peer

import (
	"fmt"
	"strings"

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

	// make sure it's a P2P addr
	parts := ma.Split(m)
	if len(parts) < 1 {
		return nil, ErrInvalidAddr
	}

	// TODO(lgierth): we shouldn't assume /p2p is the last part
	p2ppart := parts[len(parts)-1]
	if p2ppart.Protocols()[0].Code != ma.P_P2P {
		return nil, ErrInvalidAddr
	}

	// make sure the /p2p value parses as a peer.ID
	peerIdParts := strings.Split(p2ppart.String(), "/")
	peerIdStr := peerIdParts[len(peerIdParts)-1]
	id, err := IDB58Decode(peerIdStr)
	if err != nil {
		return nil, err
	}

	// we might have received just an /p2p part, which means there's no addr.
	var addrs []ma.Multiaddr
	if len(parts) > 1 {
		addrs = append(addrs, ma.Join(parts[:len(parts)-1]...))
	}

	return &AddrInfo{
		ID:    id,
		Addrs: addrs,
	}, nil
}

func AddrInfoToP2pAddrs(pi *AddrInfo) ([]ma.Multiaddr, error) {
	var addrs []ma.Multiaddr
	tpl := "/" + ma.ProtocolWithCode(ma.P_P2P).Name + "/"
	for _, addr := range pi.Addrs {
		p2paddr, err := ma.NewMultiaddr(tpl + IDB58Encode(pi.ID))
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr.Encapsulate(p2paddr))
	}
	return addrs, nil
}

func (pi *AddrInfo) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"peerID": pi.ID.Pretty(),
		"addrs":  pi.Addrs,
	}
}
