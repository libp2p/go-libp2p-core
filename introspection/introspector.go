package introspection

import (
	"github.com/libp2p/go-libp2p/core/introspection"
)

// Introspector is the interface to be satisfied by components that are capable
// of spelunking the state of the system, and representing in accordance with
// the introspection schema.
//
// It's very rare to build a custom implementation of this interface;
// it exists mostly for mocking. In most cases, you'll end up using the
// default introspector.
//
// Introspector implementations are usually injected in introspection endpoints
// to serve the data to clients, but they can also be used separately for
// embedding or testing.
//
// Experimental.
// Deprecated: use github.com/libp2p/go-libp2p/core/introspection.Introspector instead
type Introspector = introspection.Introspector
