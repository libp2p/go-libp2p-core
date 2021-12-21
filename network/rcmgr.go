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
	// GetSystem retrieves the system wide resource scope
	GetSystem() ResourceScope
	// GetDMZ retrieves the DMZ resource scope
	GetDMZ() ResourceScope
	// GetService retrieves a service-specific scope
	GetService(srv string) ServiceScope
	// GetProtocol retrieves the resource management scope for a specific protocol.
	// If there is no configured limits for a particular protocol, then the default scope is
	// returned.
	GetProtocol(protocol.ID) ProtocolScope
	// GetPeer retrieces the resource management scope for a specific peer.
	GetPeer(peer.ID) PeerScope

	// OpenConnection creates a connection scope not yet associated with any peer; the connection
	// is scoped at the DMZ.
	OpenConnection(dir Direction, usefd bool) (ConnectionScope, error)

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
	GetBuffer(size int) ([]byte, error)
	// GrowBuffer atomically grows a previous allocated buffer, reserving the appropriate memory space
	// and releasing the old buffer. The copy parameter specifies the number of bytes to copy from
	// the old buffer to the newly allocated buffer.
	GrowBuffer(buf []byte, newsize, copy int) ([]byte, error)
	// ReleaseBuffer releases a previous allocated buffer.
	ReleaseBuffer(buf []byte)

	// Stat retrieves current resource usage for the scope.
	Stat() ScopeStat
}

// TransactionalScope is a mixin interface for transactional scopes.
type TransactionalScope interface {
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

	// Protocols returns the list of protocol IDs constrained by this scope.
	Protocols() []protocol.ID
}

// PeerScope is the interface for peer resource scopes.
type PeerScope interface {
	ResourceScope

	// Peer returns the peer ID for this scope
	Peer() peer.ID

	// OpenSconnect creates a new connection scope for this peer.
	OpenConnection(dir Direction, usefd bool) (ConnectionScope, error)

	// OpenStream creates a new stream scope, with the specified protocols.
	// An unnegotiated stream will have an empty protocol list and be initially unattached to any
	// protocol scope.
	OpenStream(dir Direction, proto ...protocol.ID) (StreamScope, error)
}

// ConnectionScope is the interface for connection resource scopes.
type ConnectionScope interface {
	ResourceScope
	TransactionalScope

	// PeerScope returns the peer scope associated with this connection.
	// It reeturns nil if the connection is not yet asociated with any peer.
	PeerScope() PeerScope

	// SetPeer sets the peer for a previously unassociated connection
	SetPeer(peer.ID) error
}

// StreamScope is the interface for stream resource scopes
type StreamScope interface {
	ResourceScope
	TransactionalScope

	// ProtocolScope returns the protocol resource scope associated with this stream.
	// It returns nil if the stream is not associated with any scope.
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

// ScopeStat is a struct containing resource accounting information.
type ScopeStat struct {
	NumPeers   int
	NumConns   int
	NumStreams int

	Memory int64
}
