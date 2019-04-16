package peer

import (
	"fmt"
	"strings"

	ma "github.com/multiformats/go-multiaddr"
)

// PeerInfo is a small struct used to pass around a peer with
// a set of addresses (and later, keys?).
type Info struct {
	ID    ID
	Addrs []ma.Multiaddr
}

var _ fmt.Stringer = Info{}

func (pi Info) String() string {
	return fmt.Sprintf("{%v: %v}", pi.ID, pi.Addrs)
}

var ErrInvalidAddr = fmt.Errorf("invalid p2p multiaddr")

func InfoFromP2pAddr(m ma.Multiaddr) (*Info, error) {
	if m == nil {
		return nil, ErrInvalidAddr
	}

	// make sure it's an IPFS addr
	parts := ma.Split(m)
	if len(parts) < 1 {
		return nil, ErrInvalidAddr
	}

	// TODO(lgierth): we shouldn't assume /ipfs is the last part
	ipfspart := parts[len(parts)-1]
	if ipfspart.Protocols()[0].Code != ma.P_IPFS {
		return nil, ErrInvalidAddr
	}

	// make sure the /ipfs value parses as a peer.ID
	peerIdParts := strings.Split(ipfspart.String(), "/")
	peerIdStr := peerIdParts[len(peerIdParts)-1]
	id, err := IDB58Decode(peerIdStr)
	if err != nil {
		return nil, err
	}

	// we might have received just an /ipfs part, which means there's no addr.
	var addrs []ma.Multiaddr
	if len(parts) > 1 {
		addrs = append(addrs, ma.Join(parts[:len(parts)-1]...))
	}

	return &Info{
		ID:    id,
		Addrs: addrs,
	}, nil
}

func InfoToP2pAddrs(pi *Info) ([]ma.Multiaddr, error) {
	var addrs []ma.Multiaddr
	tpl := "/" + ma.ProtocolWithCode(ma.P_IPFS).Name + "/"
	for _, addr := range pi.Addrs {
		p2paddr, err := ma.NewMultiaddr(tpl + IDB58Encode(pi.ID))
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr.Encapsulate(p2paddr))
	}
	return addrs, nil
}

func (pi *Info) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"peerID": pi.ID.Pretty(),
		"addrs":  pi.Addrs,
	}
}
