package canonicallog

import (
	"net"
	"strings"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

var log = logging.WithSkip(logging.Logger("canonical-log"), 1)

// LogMisbehavingPeer is the canonical way to log a misbehaving peer.
// Protocols should use this to identify a misbehaving peer to allow the end
// user to easily identify these nodes across protocols and libp2p.
func LogMisbehavingPeer(p peer.ID, peerAddr multiaddr.Multiaddr, err error, msg string) {
	log.Errorf("CANONICAL_MISBEHAVING_PEER: peer=%s addr=%s err=%v msg=%s", p, peerAddr.String(), err, msg)
}

// LogMisbehavingPeer is the canonical way to log a misbehaving peer.
// Protocols should use this to identify a misbehaving peer to allow the end
// user to easily identify these nodes across protocols and libp2p.
func LogMisbehavingPeerNetAddr(p peer.ID, peerAddr net.Addr, originalErr error, msg string) {
	ipStrandPort := strings.Split(peerAddr.String(), ":")
	ip := net.ParseIP(ipStrandPort[0])
	if ip == nil {
		log.Errorf("CANONICAL_MISBEHAVING_PEER: peer=%s err=%v msg=%s", p, originalErr, msg)
		return
	}

	proto := peerAddr.Network()

	var stringBuilder strings.Builder

	if ip4 := ip.To4(); ip4 != nil {
		stringBuilder.WriteString("/ip4/")
		stringBuilder.WriteString(ip4.String())
	} else {
		stringBuilder.WriteString("/ip6/")
		stringBuilder.WriteString(ip.String())
	}
	stringBuilder.WriteString("/")

	stringBuilder.WriteString(proto)
	stringBuilder.WriteString("/")
	stringBuilder.WriteString(ipStrandPort[1])

	log.Errorf("CANONICAL_MISBEHAVING_PEER: peer=%s addr=%s err=%v msg=%s", p, stringBuilder.String(), originalErr, msg)
}
