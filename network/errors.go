package network

import (
	"errors"
	"net"
)

type TemporaryError string

func (e TemporaryError) Error() string   { return string(e) }
func (e TemporaryError) Temporary() bool { return true }
func (e TemporaryError) Timeout() bool   { return false }

var _ net.Error = TemporaryError("")

// ErrNoRemoteAddrs is returned when there are no addresses associated with a peer during a dial.
var ErrNoRemoteAddrs = errors.New("no remote addresses")

// ErrNoConn is returned when attempting to open a stream to a peer with the NoDial
// option and no usable connection is available.
var ErrNoConn = errors.New("no usable connection to peer")

// ErrTransientConn is returned when attempting to open a stream to a peer with only a transient
// connection, without specifying the UseTransient option.
var ErrTransientConn = errors.New("transient connection to peer")

// ErrResourceLimitExceeded is returned when attempting to perform an operation that would
// exceed system resource limits.
var ErrResourceLimitExceeded = TemporaryError("resource limit exceeded")

// ErrResourceScopeClosed is returned when attemptig to reserve resources in a closed resource
// scope.
var ErrResourceScopeClosed = errors.New("resource scope closed")
