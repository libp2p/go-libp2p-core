// Package transport provides the Transport interface, which represents
// the devices and network protocols used to send and receive data.
package transport

import (
	"context"
	"net"
	"time"

	"github.com/libp2p/go-libp2p-core/mux"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"

	ma "github.com/multiformats/go-multiaddr"
)

// DialTimeout is the maximum duration a Dial is allowed to take.
// This includes the time between dialing the raw network connection,
// protocol selection as well the handshake, if applicable.
var DialTimeout = 60 * time.Second

// AcceptTimeout is the maximum duration an Accept is allowed to take.
// This includes the time between accepting the raw network connection,
// protocol selection as well as the handshake, if applicable.
var AcceptTimeout = 60 * time.Second

// A CapableConn represents a connection that has offers the basic
// capabilities required by libp2p: stream multiplexing, encryption and
// peer authentication.
//
// DEPREACATED try to use QCapableConn instead, using this will result in
// wrong transport choice when trying to pick a good one.
//
// These capabilities may be natively provided by the transport, or they
// may be shimmed via the "connection upgrade" process, which converts a
// "raw" network connection into one that supports such capabilities by
// layering an encryption channel and a stream multiplexer.
//
// CapableConn provides accessors for the local and remote multiaddrs used to
// establish the connection and an accessor for the underlying Transport.
type CapableConn interface {
	mux.MuxedConn
	network.ConnSecurity
	network.ConnMultiaddrs

	// Transport returns the transport to which this connection belongs.
	Transport() Transport
}

// Quality is an estimation of how bad the connection will be.
// It must be short and fast to run.
// Lower is the returned score, better we expect the connection to be.
// Don't implement a pinging/monitoring algorithm quality must be deterministic.
//
// To see where to place your transport on the scale take a look at this (note:
// a proto with a Quality twice bigger is not twice slower or twice worst, it is
// just worst, like the place on a leaderboard):
// - TCP, 2^31, multiplexed stated continuous protocol (stream have fight for a
//   place in the underlying protocol)
// - QUIC, 2^30, async stated discontinuous protocol (stream doesn't have to
//   fight but there is time taken to open the connection and stream)
// - Circuit, router underlying connection quality + the number of hops * 8
//   (that consider that all router are equal but there is really no better than
//   pinging/monitoring to know that) + 2^16 (base circuit value).
//
// Its also not needed to follow that closely, thing can be added if a proto add
// some overhead :
// - WS tcp score + 50
// - WSS ws score + 150
//
// If conn is on a private network Quality() must divide score by 2^8 (shift
// right by 8).
// If conn is on the loopback Quality() must divide score by 2^16 (shift right
// by 16).
//
// QCapableConn embed CapableConn but with `Quality() uint32` support.
type QCapableConn interface {
	CapableConn

	// Quality returns the Quality we can expect from the connection to this peer.
	// That must be deterministic and fast.
	Quality() uint32
}

// Score is used by transport to returns expectation about connection
type Score struct {
	// Score of the future connection, must not change once the connection is
	// actualy created
	Quality uint32
	// True if quality was actualy calculable (it may only be calculable with an
	// open connection).
	// Note: if IsQuality == true swarm will be able to cancel the dial if a
	// better transport is found.
	IsQuality bool
	// Number of file descriptor expected to be open.
	Fd uint8
}

// Transport represents any device by which you can connect to and accept
// connections from other peers.
//
// The Transport interface allows you to open connections to other peers
// by dialing them, and also lets you listen for incoming connections.
//
// Connections returned by Dial and passed into Listeners are of type
// CapableConn, which means that they have been upgraded to support
// stream multiplexing and connection security (encryption and authentication).
//
// For a conceptual overview, see https://docs.libp2p.io/concepts/transport/
//
// DEPREACATED try to use QTransport instead, using this will result in
// wrong transport choice when trying to pick a good one.
type Transport interface {
	BaseTransport

	// Dial dials a remote peer. It should try to reuse local listener
	// addresses if possible but it may choose not to.
	Dial(ctx context.Context, raddr ma.Multiaddr, p peer.ID) (CapableConn, error)
}

type QTransport interface {
	BaseTransport

	// Dial dials a remote peer. It should try to reuse local listener
	// addresses if possible but it may choose not to.
	Dial(ctx context.Context, raddr ma.Multiaddr, p peer.ID) (QCapableConn, error)

	// Score returns the Quality we can expect from the connection to this peer.
	// That must be deterministic and fast.
	Score(raddr ma.Multiaddr, p peer.ID) (Score, error)
}

type BaseTransport interface {
	// CanDial returns true if this transport knows how to dial the given
	// multiaddr.
	//
	// Returning true does not guarantee that dialing this multiaddr will
	// succeed. This function should *only* be used to preemptively filter
	// out addresses that we can't dial.
	CanDial(addr ma.Multiaddr) bool

	// Listen listens on the passed multiaddr.
	Listen(laddr ma.Multiaddr) (Listener, error)

	// Protocol returns the set of protocols handled by this transport.
	//
	// See the Network interface for an explanation of how this is used.
	Protocols() []int

	// Proxy returns true if this is a proxy transport.
	//
	// See the Network interface for an explanation of how this is used.
	// TODO: Make this a part of the go-multiaddr protocol instead?
	Proxy() bool
}

// Listener is an interface closely resembling the net.Listener interface. The
// only real difference is that Accept() returns Conn's of the type in this
// package, and also exposes a Multiaddr method as opposed to a regular Addr
// method
type Listener interface {
	Accept() (CapableConn, error)
	Close() error
	Addr() net.Addr
	Multiaddr() ma.Multiaddr
}

// Network is an inet.Network with methods for managing transports.
type TransportNetwork interface {
	network.Network

	// AddTransport adds a transport to this Network.
	//
	// When dialing, this Network will iterate over the protocols in the
	// remote multiaddr and pick the first protocol registered with a proxy
	// transport, if any. Otherwise, it'll pick the transport registered to
	// handle the last protocol in the multiaddr.
	//
	// When listening, this Network will iterate over the protocols in the
	// local multiaddr and pick the *last* protocol registered with a proxy
	// transport, if any. Otherwise, it'll pick the transport registered to
	// handle the last protocol in the multiaddr.
	//
	// Even if BaseTransport are accepted its needed to be able to cast to
	// Transport or QTransport.
	AddTransport(t BaseTransport) error
}
