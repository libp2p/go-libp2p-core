package introspection

// ProvidersTree is a struct which provider modules use to register their interest in providing
// data to be recorded in a time slice of the running system.
type ProvidersTree struct {
	Runtime *RuntimeProviders
	Conn    *ConnProviders
	Stream  *StreamProviders
}

// synthetic used types to contextualise inputs/outputs and enforce correctness.
type (
	// ConnID represents a connection ID string
	ConnID string
	// StreamID represents a connection ID string
	StreamID string
)

// Enums.
type (
	// StreamListQueryType is an Enum to represent the types of queries that can be made when looking up streams
	StreamListQueryType int
)

const (
	// StreamListQueryTypeAll is an enum to represent a query to fetch all streams
	StreamListQueryTypeAll StreamListQueryType = iota
	// StreamListQueryTypeConn is an enum to represent a query to fetch all streams for connections
	StreamListQueryTypeConn
)

// ConnProviders contains provider functions that deal with connections. See
// godoc on StreamProviders for an explanation of what List and Get do.
type ConnProviders struct {
	List func() ([]*Connection, error)
	Get  func([]ConnID) ([]*Connection, error)
}

// StreamListQuery contrextualises query params for getting streams
type StreamListQuery struct {
	Type   StreamListQueryType
	ConnID ConnID
}

// StreamProviders contains provider functions that deal with streams.
type StreamProviders struct {
	// List is a provider function that returns a shallow list of streams
	// matching a query.
	List func(query StreamListQuery) ([]*Stream, error)

	// Get returns populated Stream objects for the requested streams.
	Get func([]StreamID) ([]*Stream, error)
}

// RuntimeProviders contains provider functions that deal with the runtime.
type RuntimeProviders struct {
	Get func() (*Runtime, error)
}
