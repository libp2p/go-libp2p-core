// Deprecated: This package has moved into go-libp2p as a sub-package: github.com/libp2p/go-libp2p/core.
//
// Package core provides convenient access to foundational, central go-libp2p primitives via type aliases.
package core

import (
	"github.com/libp2p/go-libp2p/core"
)

// Multiaddr aliases the Multiaddr type from github.com/multiformats/go-multiaddr.
//
// Refer to the docs on that type for more info.
// Deprecated: use github.com/libp2p/go-libp2p/core.Multiaddr instead
type Multiaddr = core.Multiaddr

// PeerID aliases peer.ID.
//
// Refer to the docs on that type for more info.
// Deprecated: use github.com/libp2p/go-libp2p/core.PeerID instead
type PeerID = core.PeerID

// ProtocolID aliases protocol.ID.
//
// Refer to the docs on that type for more info.
// Deprecated: use github.com/libp2p/go-libp2p/core.ProtocolID instead
type ProtocolID = core.ProtocolID

// PeerAddrInfo aliases peer.AddrInfo.
//
// Refer to the docs on that type for more info.
// Deprecated: use github.com/libp2p/go-libp2p/core.PeerAddrInfo instead
type PeerAddrInfo = core.PeerAddrInfo

// Host aliases host.Host.
//
// Refer to the docs on that type for more info.
// Deprecated: use github.com/libp2p/go-libp2p/core.Host instead
type Host = core.Host

// Network aliases network.Network.
//
// Refer to the docs on that type for more info.
// Deprecated: use github.com/libp2p/go-libp2p/core.Network instead
type Network = core.Network

// Conn aliases network.Conn.
//
// Refer to the docs on that type for more info.
// Deprecated: use github.com/libp2p/go-libp2p/core.Conn instead
type Conn = core.Conn

// Stream aliases network.Stream.
//
// Refer to the docs on that type for more info.
// Deprecated: use github.com/libp2p/go-libp2p/core.Stream instead
type Stream = core.Stream
