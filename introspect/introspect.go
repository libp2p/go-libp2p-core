package introspect

import introspect_pb "github.com/libp2p/go-libp2p-core/introspect/pb"

// ProtoVersion is the current version of the introspection protocol.
const ProtoVersion uint32 = 1

// EXPERIMENTAL. Introspector allows other sub-systems/modules to register
// metrics/data providers AND also enables clients to fetch the current state of
// the system.
type Introspector interface {

	// EXPERIMENTAL. RegisterDataProviders allows sub-systems/modules to
	// register callbacks that supply introspection data.
	RegisterDataProviders(p *DataProviders) error

	// EXPERIMENTAL. FetchFullState returns the full state of the system, by
	// calling all known data providers and returning a merged cross-cut of the
	// running system.
	FetchFullState() (*introspect_pb.State, error)

	// EXPERIMENTAL. ListenAddrs returns the addresses on which the
	// introspection server endpoint is listening for clients.
	ListenAddrs() []string
}
