package introspection

import (
	"io"

	pb "github.com/libp2p/go-libp2p-core/introspection/pb"
)

// ProtoVersion is the current version of the introspection protocol.
const ProtoVersion uint32 = 1

// ProtoVersionPb is the proto representation of the current introspection protocol.
var ProtoVersionPb = &pb.Version{Version: ProtoVersion}

// EXPERIMENTAL. Introspector allows introspection endpoints to fetch the current state of
// the system.
//
// Introspector implementations are usually injected in introspection endpoints
// (e.g. the default WebSocket endpoint) to serve the data to clients, but they
// can also be used separately for embedding or testing.
type Introspector interface {
	io.Closer

	// FetchRuntime returns the runtime information of the system.
	// Experimental.
	FetchRuntime() (*pb.Runtime, error)

	// FetchFullState returns the full state cross-cut of the running system.
	// Experimental.
	FetchFullState() (*pb.State, error)

	// EventChan returns the channel where all eventbus events are dumped,
	// decorated with their corresponding event metadata, ready to send over
	// the wire.
	// Experimental.
	EventChan() <-chan *pb.Event

	// EventMetadata returns the metadata of all events known to the
	// Introspector.
	// Experimental
	EventMetadata() []*pb.EventType
}
