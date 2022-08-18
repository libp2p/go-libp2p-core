// Deprecated: This package has moved into go-libp2p as a sub-package: github.com/libp2p/go-libp2p/core/transport.
//
// Package transport provides the Transport interface, which represents
// the devices and network protocols used to send and receive data.
package transport

import (
	"github.com/libp2p/go-libp2p/core/transport"
)

// A CapableConn represents a connection that has offers the basic
// capabilities required by libp2p: stream multiplexing, encryption and
// peer authentication.
//
// These capabilities may be natively provided by the transport, or they
// may be shimmed via the "connection upgrade" process, which converts a
// "raw" network connection into one that supports such capabilities by
// layering an encryption channel and a stream multiplexer.
//
// CapableConn provides accessors for the local and remote multiaddrs used to
// establish the connection and an accessor for the underlying Transport.
// Deprecated: use github.com/libp2p/go-libp2p/core/transport.CapableConn instead
type CapableConn = transport.CapableConn

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
// If a transport implements `io.Closer` (optional), libp2p will call `Close` on
// shutdown. NOTE: `Dial` and `Listen` may be called after or concurrently with
// `Close`.
//
// For a conceptual overview, see https://docs.libp2p.io/concepts/transport/
// Deprecated: use github.com/libp2p/go-libp2p/core/transport.Transport instead
type Transport = transport.Transport

// Listener is an interface closely resembling the net.Listener interface. The
// only real difference is that Accept() returns Conn's of the type in this
// package, and also exposes a Multiaddr method as opposed to a regular Addr
// method
// Deprecated: use github.com/libp2p/go-libp2p/core/transport.Listener instead
type Listener = transport.Listener

// TransportNetwork is an inet.Network with methods for managing transports.
// Deprecated: use github.com/libp2p/go-libp2p/core/transport.TransportNetwork instead
type TransportNetwork = transport.TransportNetwork

// Upgrader is a multistream upgrader that can upgrade an underlying connection
// to a full transport connection (secure and multiplexed).
// Deprecated: use github.com/libp2p/go-libp2p/core/transport.Upgrader instead
type Upgrader = transport.Upgrader
