// Deprecated: This package has moved into go-libp2p as a sub-package: github.com/libp2p/go-libp2p/core/pnet.
//
// Package pnet provides interfaces for private networking in libp2p.
package pnet

import "github.com/libp2p/go-libp2p/core/pnet"

// A PSK enables private network implementation to be transparent in libp2p.
// It is used to ensure that peers can only establish connections to other peers
// that are using the same PSK.
// Deprecated: use github.com/libp2p/go-libp2p/core/pnet.PSK instead
type PSK = pnet.PSK
