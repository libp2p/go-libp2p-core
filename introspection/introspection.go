package introspection

// ProtoVersion is the current version of the Proto
const ProtoVersion uint32 = 1

// IntrospectorRegistry allows other sub-systems/modules to register metrics/data providers.
type IntrospectorRegistry interface {
	// RegisterProviders allows a subsystem to register itself as a provider of metrics.
	RegisterProviders(p *ProvidersTree) error
}

// Introspector allows other sub-systems/modules to register metrics/data providers AND also
// enables clients to fetch the current state of the system.
type Introspector interface {
	IntrospectorRegistry

	// FetchCurrentState fetches the current state of the sub-systems by calling the providers registered by them on the registry.
	FetchCurrentState() (*State, error)
}
