package introspection

// synthetic used types to contextualise inputs/outputs and enforce correctness.
type (
	// ConnID represents a connection ID string
	ConnID string
	// StreamID represents a connection ID string
	StreamID string
)

// ProvidersMap is a struct which provider modules use to register their interest in providing
// data to be recorded in a time slice of the running system.
type ProvidersMap struct {
	Runtime    func() (*Runtime, error)
	Connection func(ConnectionQueryInput) ([]*Connection, error)
	Stream     func(StreamQueryInput) ([]*Stream, error)
}

type (
	// StreamListQueryType is an Enum to represent the types of queries that can be made when looking up streams
	StreamListQueryType int
)

const (
	// StreamListQueryTypeAll is an enum to represent a query to fetch all streams
	StreamListQueryTypeAll StreamListQueryType = iota
	// StreamListQueryTypeConn is an enum to represent a query to fetch all streams a given connection
	StreamListQueryTypeConn
	// StreamListQueryTypeIds represents a query to fetch streams with the given Ids
	StreamListQueryTypeIds
)

// StreamQueryInput determines the input for the stream query
type StreamQueryInput struct {
	Type      StreamListQueryType
	ConnID    ConnID
	StreamIds []StreamID
}

type (
	//  ConnListQueryType is an Enum to represent the types of queries that can be made when looking up streams
	ConnListQueryType int
)

const (
	// ConnListQueryTypeAll represents a query to fetch all connections
	ConnListQueryTypeAll ConnListQueryType = iota
	// ConnListQueryTypeForIds represents a query to fetch connections for a given set of Ids
	ConnListQueryTypeForIds
)

// ConnectionQueryInput determines the input for the connections query
type ConnectionQueryInput struct {
	Type    ConnListQueryType
	ConnIDs []ConnID
}
