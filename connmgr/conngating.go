package connmgr

import (
	"github.com/libp2p/go-libp2p-core/control"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"

	ma "github.com/multiformats/go-multiaddr"
)

// ConnectionGater can be implemented by a type that supports active
// inbound or outbound connection gating.
//
// A ConnectionGater will be consulted during different states in the life-cycle of a connection and
// the specific gating function that will be called depends on the life-cycle state of the
// connection. Hence, it is important to implement this interface keeping in mind the specific life-cycle state
// at which you'd like to gate/block the connection.
//
// `InterceptDial` and `InterceptPeerDial` are called when we try an outbound dial.
//
// `InterceptAccept` is called after we've accepted an inbound connection from a socket but before we
// begin upgrading it.
//
// `InterceptSecured` is called for inbound and outbound connections after we've negotiated
// the security protocol to use for the connection.
//
// `InterceptUpgraded` is called for inbound and outbound connections after we've finished upgrading
// a connection to have both security and stream multiplexing.
//
// This feature can be used to implement *strict* connection management
// behaviours, such as hard limiting of connections once a max count has been
// reached.
//
// If you'd like to send a disconnect control message to the remote peer for a gated inbound connection,
// ONLY `InterceptUpgraded` should reject the connection (with an appropriate disconenct reason).
// All other methods should allow the connection as we can ONLY open control streams
// for upgraded connections.
// Note: There's no point in sending disconnect control messages for outbound connections, so we might
// as well close them as early in the cycle as possible.
type ConnectionGater interface {
	// InterceptDial tests whether we're permitted to dial the specified multiaddr.
	// Insofar filter.Filters is concerned, this would map to its AddrBlock method,
	// with the inverse condition.
	// This is to be called by the network/swarm when dialling.
	InterceptDial(ma.Multiaddr) (allow bool)

	// InterceptPeerDial tests whether we're permitted to Dial the specified peer.
	// This is to be called by the network/swarm when dialling.
	InterceptPeerDial(p peer.ID) (allow bool)

	// InterceptAccept tests whether an incipient inbound connection is allowed.
	// network.ConnMultiaddrs is what we pass to the upgrader.
	// This is intended to be called by the upgrader, or by the transport
	// directly (e.g. QUIC, Bluetooth), straight after it's accepted a connection
	// from its socket.
	InterceptAccept(network.ConnMultiaddrs) (allow bool)

	// InterceptSecured tests whether a given connection, now authenticated,
	// is allowed.
	// This is intended to be called by the upgrader, after it has negotiated crypto,
	// and before it negotiates the muxer, or by the directly by the transport,
	// at the exact same checkpoint.
	InterceptSecured(network.Direction, peer.ID, network.ConnMultiaddrs) (allow bool)

	// InterceptUpgraded tests whether a fully capable connection is allowed.
	// At this point, we have a multiplexer, so the middleware can
	// return a DisconnectReason and the swarm would use the control stream to convey
	// it to the peer.
	InterceptUpgraded(transport.CapableConn) (allow bool, reason control.DisconnectReason)
}
