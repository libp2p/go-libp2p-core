package introspection

import "github.com/libp2p/go-libp2p/core/introspection"

// Endpoint is the interface to be implemented by introspection endpoints.
//
// An introspection endpoint makes introspection data accessible to external
// consumers, over, for example, WebSockets, or TCP, or libp2p itself.
//
// Experimental.
// Deprecated: use github.com/libp2p/go-libp2p/core/introspection.Endpoint instead
type Endpoint = introspection.Endpoint

// Session represents an introspection session.
// Deprecated: use github.com/libp2p/go-libp2p/core/introspection.Session instead
type Session = introspection.Session
