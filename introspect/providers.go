package introspect

// ConnID represents a connection ID string
type ConnID string

// QueryOutputType determines the way a output needs to be represented in a query result
type QueryOutputType int

const (
	// QueryOutputTypeFull dictates that we need to resolve the whole object in the query output
	QueryOutputTypeFull QueryOutputType = iota
	// QueryOutputTypeIds dictates that we need to resolve only the identifiers of the object in the query output
	QueryOutputTypeIds
)

// ProvidersMap is a struct which provider modules use to register their interest in providing
// data to be recorded in a time slice of the running system.
type ProvidersMap struct {
	Runtime    func() (*Runtime, error)
	Connection func(ConnectionQueryInput) ([]*Connection, error)
}

// ConnListQueryType is an Enum to represent the types of queries that can be made when looking up streams
type ConnListQueryType int

const (
	// ConnListQueryTypeAll represents a query to fetch all connections
	ConnListQueryTypeAll ConnListQueryType = iota
	// ConnListQueryTypeForIds represents a query to fetch connections for a given set of Ids
	ConnListQueryTypeForIds
)

// ConnectionQueryInput determines the input for the connections query
type ConnectionQueryInput struct {
	Type             ConnListQueryType
	StreamOutputType QueryOutputType
	ConnIDs          []ConnID
}
