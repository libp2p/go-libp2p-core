package network

import (
	"io"
	"time"

	ic "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"

	ma "github.com/multiformats/go-multiaddr"
)

// Conn is a connection to a remote peer. It multiplexes streams.
// Usually there is no need to use a Conn directly, but it may
// be useful to get information about the peer on the other side:
//  stream.Conn().RemotePeer()
type Conn interface {
	io.Closer

	ConnSecurity
	ConnMultiaddrs

	// NewStream constructs a new Stream over this conn.
	NewStream() (Stream, error)

	// GetStreams returns all open streams over this conn.
	GetStreams() []Stream

	// Stat stores metadata pertaining to this conn.
	Stat() Stat

	// OnBetter is an callback when a better connection is found.
	// OnBetter is threadsafe, it can be called even once the event raised and the
	// callback we be yield.
	OnBetter(OnBetterHandler)
}

// OnBetterHandler are args to pass to Conn.OnBetter The time is the deadline
// before a hard close (fixed time not duration).
type OnBetterHandler func(time.Time)

// ConnSecurity is the interface that one can mix into a connection interface to
// give it the security methods.
type ConnSecurity interface {
	// LocalPeer returns our peer ID
	LocalPeer() peer.ID

	// LocalPrivateKey returns our private key
	LocalPrivateKey() ic.PrivKey

	// RemotePeer returns the peer ID of the remote peer.
	RemotePeer() peer.ID

	// RemotePublicKey returns the public key of the remote peer.
	RemotePublicKey() ic.PubKey
}

// ConnMultiaddrs is an interface mixin for connection types that provide multiaddr
// addresses for the endpoints.
type ConnMultiaddrs interface {
	// LocalMultiaddr returns the local Multiaddr associated
	// with this connection
	LocalMultiaddr() ma.Multiaddr

	// RemoteMultiaddr returns the remote Multiaddr associated
	// with this connection
	RemoteMultiaddr() ma.Multiaddr
}
