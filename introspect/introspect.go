package introspect

// ProtoVersion is the current version of the Proto
const ProtoVersion uint32 = 1

// Introspector allows other sub-systems/modules to register metrics/data providers AND also
// enables clients to fetch the current state of the system.
type Introspector interface {
	// RegisterProviders allows sub-systems/modules to register themselves as data/metrics providers
	RegisterProviders(p *ProvidersMap) error

	// FetchCurrentState fetches the current state of the sub-systems by calling the providers registered by them on the registry.
	FetchCurrentState() (*State, error)

	// ListenAddress returns the address on which the introspection service will be available
	ListenAddress() string
}
