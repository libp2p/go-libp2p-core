package canonicallog

import (
	"net"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

var log = logging.WithSkip(logging.Logger("canonical-log"), 1)

// LogMisbehavingPeer is the canonical way to log a misbehaving peer.
// Protocols should use this to identify a misbehaving peer to allow the end
// user to easily identify these nodes across protocols and libp2p.
func LogMisbehavingPeer(p peer.ID, peerAddr multiaddr.Multiaddr, component string, err error, msg string) {
	log.Errorf("CANONICAL_MISBEHAVING_PEER: peer=%s addr=%s component=%s err=%v msg=%s", p, peerAddr.String(), component, err, msg)
}

// LogMisbehavingPeer is the canonical way to log a misbehaving peer.
// Protocols should use this to identify a misbehaving peer to allow the end
// user to easily identify these nodes across protocols and libp2p.
func LogMisbehavingPeerNetAddr(p peer.ID, peerAddr net.Addr, component string, originalErr error, msg string) {
	ma, err := manet.FromNetAddr(peerAddr)
	if err != nil {
		log.Errorf("CANONICAL_MISBEHAVING_PEER: peer=%s netAddr=%s component=%s err=%v msg=%s", p, peerAddr.String(), component, originalErr, msg)
	}

	log.Errorf("CANONICAL_MISBEHAVING_PEER: peer=%s addr=%s component=%s err=%v msg=%s", p, ma, component, originalErr, msg)
}
