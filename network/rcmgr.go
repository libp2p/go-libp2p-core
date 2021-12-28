package network

import (
	"errors"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

var (
	ErrResourceLimitExceeded = errors.New("resource limit exceeded")
	ErrResourceScopeClosed   = errors.New("resource scope closed")
)

// ResourceManager is the interface to the network resource management subsystem
type ResourceManager interface {
	// ViewSystem views the system wide resource scope
	ViewSystem(func(ResourceScope) error) error
	// ViewTransient views the transient (DMZ) resource scope
	ViewTransient(func(ResourceScope) error) error
	// ViewService retrieves a service-specific scope
	ViewService(string, func(ServiceScope) error) error
	// ViewProtocol views the resource management scope for a specific protocol.
	ViewProtocol(protocol.ID, func(ProtocolScope) error) error
	// ViewPeer views the resource management scope for a specific peer.
	ViewPeer(peer.ID, func(PeerScope) error) error

	// OpenConnection creates a new connection scope not yet associated with any peer; the connection
	// is scoped at the transient scope.
	OpenConnection(dir Direction, usefd bool) (ConnectionScope, error)

	// OpenStream creates a new stream scope, initially unnegotiated.
	// An unnegotiated stream will be initially unattached to any protocol scope
	// and constrained by the transient scope.
	OpenStream(p peer.ID, dir Direction) (StreamScope, error)

	// Close closes the resource manager
	Close() error
}

// ResourceScope is the interface for all scopes.
type ResourceScope interface {
	// ReserveMemory reserves memory/buffer space in the scope.
	ReserveMemory(size int) error
	// ReleaseMemory explicitly releases memory previously reserved with ReserveMemory
	ReleaseMemory(size int)

	// GetBuffer reserves memory and allocates a buffer through the buffer pool.
	GetBuffer(size int) (Buffer, error)

	// Stat retrieves current resource usage for the scope.
	Stat() ScopeStat

	// BeginTxn creates a new transactional scope rooted at this scope
	BeginTxn() (TransactionalScope, error)
}

// Buffer is a wrapper for a buffer managed by a resource scope
type Buffer interface {
	// Data returns the byte slice for this buffer
	Data() []byte
	// Grow grows the buffer to a new size
	Grow(newsize int) error
	// Release releases the buffer
	Release()
}

// TransactionalScope is a ResourceScope with transactional semantics.
type TransactionalScope interface {
	ResourceScope
	// Done ends the transaction scope and releases associated resources.
	Done()
}

// ServiceScope is the interface for service resource scopes
type ServiceScope interface {
	ResourceScope

	// Name returns the name of this service
	Name() string
}

// ProtocolScope is the interface for protocol resource scopes.
type ProtocolScope interface {
	ResourceScope

	// Protocol returns the list of protocol IDs constrained by this scope.
	Protocol() protocol.ID
}

// PeerScope is the interface for peer resource scopes.
type PeerScope interface {
	ResourceScope

	// Peer returns the peer ID for this scope
	Peer() peer.ID
}

// ConnectionScope is the interface for connection resource scopes.
type ConnectionScope interface {
	TransactionalScope

	// PeerScope returns the peer scope associated with this connection.
	// It reeturns nil if the connection is not yet asociated with any peer.
	PeerScope() PeerScope

	// SetPeer sets the peer for a previously unassociated connection
	SetPeer(peer.ID) error
}

// UserConnectionScope is the user view of a ConnectionScope
type UserConnectionScope interface {
	ResourceScope
}

// StreamScope is the interface for stream resource scopes.
type StreamScope interface {
	TransactionalScope

	// ProtocolScope returns the protocol resource scope associated with this stream.
	// It returns nil if the stream is not associated with any protocol scope.
	ProtocolScope() ProtocolScope
	// SetProtocol sets the protocol for a previously unnegotiated stream
	SetProtocol(proto protocol.ID) error

	// ServiceScope returns the service owning the stream, if any.
	ServiceScope() ServiceScope
	// SetService sets the service owning this stream
	SetService(srv string) error

	// PeerScope returns the peer resource scope associated with this stream.
	PeerScope() PeerScope
}

// UserStreamScope is the user view of a StreamScope
type UserStreamScope interface {
	ResourceScope

	// SetService sets the service owning this stream
	SetService(srv string) error
}

// ScopeStat is a struct containing resource accounting information.
type ScopeStat struct {
	NumStreamsInbound  int
	NumStreamsOutbound int
	NumConnsInbound    int
	NumConnsOutbound   int
	NumFD              int

	Memory int64
}
