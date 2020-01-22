package introspection

// Introspector interface should be implemented by a type that wants to allows other sub-systems/modules to register
// metrics providers. It can then get a whole picture of it's subsystems by calling the providers of each subsystem & stitching all the data together.
type Introspector interface {
	// RegisterProviders allows a subsystem to register itself as a provider of metrics
	RegisterProviders(p *ProvidersTree) error
}
