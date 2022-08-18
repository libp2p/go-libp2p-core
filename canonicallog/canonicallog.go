// Deprecated: This package has moved into go-libp2p as a sub-package: github.com/libp2p/go-libp2p/core/canonicallog.
package canonicallog

import (
	"net"

	"github.com/libp2p/go-libp2p/core/canonicallog"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/multiformats/go-multiaddr"
)

// LogMisbehavingPeer is the canonical way to log a misbehaving peer.
// Protocols should use this to identify a misbehaving peer to allow the end
// user to easily identify these nodes across protocols and libp2p.
// Deprecated: use github.com/libp2p/go-libp2p/core/canonicallog.LogMisbehavingPeer instead
func LogMisbehavingPeer(p peer.ID, peerAddr multiaddr.Multiaddr, component string, err error, msg string) {
	canonicallog.LogMisbehavingPeer(p, peerAddr, component, err, msg)
}

// LogMisbehavingPeerNetAddr is the canonical way to log a misbehaving peer.
// Protocols should use this to identify a misbehaving peer to allow the end
// user to easily identify these nodes across protocols and libp2p.
// Deprecated: use github.com/libp2p/go-libp2p/core/canonicallog.LogMisbehavingPeerNetAddr instead
func LogMisbehavingPeerNetAddr(p peer.ID, peerAddr net.Addr, component string, originalErr error, msg string) {
	canonicallog.LogMisbehavingPeerNetAddr(p, peerAddr, component, originalErr, msg)
}

// LogPeerStatus logs any useful information about a peer. It takes in a sample
// rate and will only log one in every sampleRate messages (randomly). This is
// useful in surfacing events that are normal in isolation, but may be abnormal
// in large quantities. For example, a successful connection from an IP address
// is normal. 10,000 connections from that same IP address is not normal. libp2p
// itself does nothing besides emitting this log. Hook this up to another tool
// like fail2ban to action on the log.
// Deprecated: use github.com/libp2p/go-libp2p/core/canonicallog.LogPeerStatus instead
func LogPeerStatus(sampleRate int, p peer.ID, peerAddr multiaddr.Multiaddr, keyVals ...string) {
	canonicallog.LogPeerStatus(sampleRate, p, peerAddr, keyVals...)
}
