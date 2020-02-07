package introspection

import introspection_pb "github.com/libp2p/go-libp2p-core/introspection/pb"

type (
	// QueryOutput determines the output form of a query result.
	QueryOutput int

	// ConnectionID represents a connection ID.
	ConnectionID string

	// StreamID represents a stream ID.
	StreamID string
)

const (
	// QueryOutputFull dictates that we need to resolve the whole object in the
	// query output.
	QueryOutputFull QueryOutput = iota

	// QueryOutputList dictates that we need to resolve only the identifiers of
	// the object in the query output.
	QueryOutputList
)

// EXPERIMENTAL. DataProviders enumerates the functions that resolve each entity
// type. It is used by go-libp2p modules to register callback functions capable
// of processing entity queries.
type DataProviders struct {
	// Runtime is the provider function that returns system runtime information.
	Runtime func() (*introspection_pb.Runtime, error)

	// Connection is the provider that is called when information about
	// Connections is required.
	Connection func(ConnectionQueryParams) ([]*introspection_pb.Connection, error)

	// Stream is the provider that is called when information about Streams is
	// required.
	Stream func(StreamQueryParams) (*introspection_pb.StreamList, error)

	// Traffic is the provider that is called when information about network
	// statistics is required.
	Traffic func() (*introspection_pb.Traffic, error)
}

type ConnectionQueryParams struct {
	Output  QueryOutput
	Include []ConnectionID
}

type StreamQueryParams struct {
	Output  QueryOutput
	Include []StreamID
}
