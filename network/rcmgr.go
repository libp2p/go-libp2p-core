package network

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
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
	OpenConnection(dir Direction, usefd bool) (ConnectionManagementScope, error)

	// OpenStream creates a new stream scope, initially unnegotiated.
	// An unnegotiated stream will be initially unattached to any protocol scope
	// and constrained by the transient scope.
	OpenStream(p peer.ID, dir Direction) (StreamManagementScope, error)

	// Close closes the resource manager
	Close() error
}

// ResourceScope is the interface for all scopes.
type ResourceScope interface {
	// ReserveMemory reserves memory/buffer space in the scope.
	ReserveMemory(size int) error
	// ReleaseMemory explicitly releases memory previously reserved with ReserveMemory
	ReleaseMemory(size int)

	// Stat retrieves current resource usage for the scope.
	Stat() ScopeStat

	// BeginTransaction creates a new transactional scope rooted at this scope
	BeginTransaction() (TransactionalScope, error)
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

	// Protocol returns the protocol for this scope
	Protocol() protocol.ID
}

// PeerScope is the interface for peer resource scopes.
type PeerScope interface {
	ResourceScope

	// Peer returns the peer ID for this scope
	Peer() peer.ID
}

// ConnectionManagementScope is the interface for connection resource scopes.
type ConnectionManagementScope interface {
	TransactionalScope

	// PeerScope returns the peer scope associated with this connection.
	// It returns nil if the connection is not yet asociated with any peer.
	PeerScope() PeerScope

	// SetPeer sets the peer for a previously unassociated connection
	SetPeer(peer.ID) error
}

// ConnectionScope is the user view of a connection scope
type ConnectionScope interface {
	ResourceScope
}

// StreamManagementScope is the interface for stream resource scopes.
type StreamManagementScope interface {
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
type StreamScope interface {
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

// NullResourceManager is a stub for tests and initialization of default values
var NullResourceManager ResourceManager = &nullResourceManager{}

type nullResourceManager struct{}
type nullScope struct{}

var _ ResourceScope = (*nullScope)(nil)
var _ TransactionalScope = (*nullScope)(nil)
var _ ServiceScope = (*nullScope)(nil)
var _ ProtocolScope = (*nullScope)(nil)
var _ PeerScope = (*nullScope)(nil)
var _ ConnectionManagementScope = (*nullScope)(nil)
var _ ConnectionScope = (*nullScope)(nil)
var _ StreamManagementScope = (*nullScope)(nil)
var _ StreamScope = (*nullScope)(nil)

var nullScopeObj = &nullScope{}

func (n *nullResourceManager) ViewSystem(f func(ResourceScope) error) error {
	return f(nullScopeObj)
}
func (n *nullResourceManager) ViewTransient(f func(ResourceScope) error) error {
	return f(nullScopeObj)
}
func (n *nullResourceManager) ViewService(svc string, f func(ServiceScope) error) error {
	return f(nullScopeObj)
}
func (n *nullResourceManager) ViewProtocol(p protocol.ID, f func(ProtocolScope) error) error {
	return f(nullScopeObj)
}
func (n *nullResourceManager) ViewPeer(p peer.ID, f func(PeerScope) error) error {
	return f(nullScopeObj)
}
func (n *nullResourceManager) OpenConnection(dir Direction, usefd bool) (ConnectionManagementScope, error) {
	return nullScopeObj, nil
}
func (n *nullResourceManager) OpenStream(p peer.ID, dir Direction) (StreamManagementScope, error) {
	return nullScopeObj, nil
}
func (n *nullResourceManager) Close() error {
	return nil
}

func (n *nullScope) ReserveMemory(size int) error                  { return nil }
func (n *nullScope) ReleaseMemory(size int)                        {}
func (n *nullScope) Stat() ScopeStat                               { return ScopeStat{} }
func (n *nullScope) BeginTransaction() (TransactionalScope, error) { return nullScopeObj, nil }
func (n *nullScope) Done()                                         {}
func (n *nullScope) Name() string                                  { return "" }
func (n *nullScope) Protocol() protocol.ID                         { return "" }
func (n *nullScope) Peer() peer.ID                                 { return "" }
func (n *nullScope) PeerScope() PeerScope                          { return nullScopeObj }
func (n *nullScope) SetPeer(peer.ID) error                         { return nil }
func (n *nullScope) ProtocolScope() ProtocolScope                  { return nullScopeObj }
func (n *nullScope) SetProtocol(proto protocol.ID) error           { return nil }
func (n *nullScope) ServiceScope() ServiceScope                    { return nullScopeObj }
func (n *nullScope) SetService(srv string) error                   { return nil }
