// Deprecated: This package has moved into go-libp2p as a sub-package: github.com/libp2p/go-libp2p/core/sec.
//
// Package sec provides secure connection and transport interfaces for libp2p.
package sec

import (
	"github.com/libp2p/go-libp2p/core/sec"
)

// SecureConn is an authenticated, encrypted connection.
// Deprecated: use github.com/libp2p/go-libp2p/core/sec.SecureConn instead
type SecureConn = sec.SecureConn

// A SecureTransport turns inbound and outbound unauthenticated,
// plain-text, native connections into authenticated, encrypted connections.
// Deprecated: use github.com/libp2p/go-libp2p/core/sec.SecureTransport instead
type SecureTransport = sec.SecureTransport

// A SecureMuxer is a wrapper around SecureTransport which can select security protocols
// and open outbound connections with simultaneous open.
// Deprecated: use github.com/libp2p/go-libp2p/core/sec.SecureMuxer instead
type SecureMuxer = sec.SecureMuxer
