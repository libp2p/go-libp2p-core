package mux

import (
	"errors"
	"io"
	"net"
	"time"
)

// ErrReset is returned when reading or writing on a reset stream.
var ErrReset = errors.New("stream reset")

// Stream is a bidirectional io pipe within a connection.
type MuxedStream interface {
	io.Reader
	io.Writer

	// Close closes the stream for writing. Reading will still work (that
	// is, the remote side can still write).
	io.Closer

	// Reset closes both ends of the stream. Use this to tell the remote
	// side to hang up and go away.
	Reset() error

	SetDeadline(time.Time) error
	SetReadDeadline(time.Time) error
	SetWriteDeadline(time.Time) error
}

// NoOpHandler do nothing. Resets streams as soon as they are opened.
var NoOpHandler = func(s MuxedStream) { s.Reset() }

// MuxedConn is a stream-multiplexing connection to a remote peer.
type MuxedConn interface {
	// Close closes the stream muxer and the the underlying net.Conn.
	io.Closer

	// IsClosed returns whether a connection is fully closed, so it can
	// be garbage collected.
	IsClosed() bool

	// OpenStream creates a new stream.
	OpenStream() (MuxedStream, error)

	// AcceptStream accepts a stream opened by the other side.
	AcceptStream() (MuxedStream, error)
}

// Transport constructs go-stream-muxer compatible connections.
type Multiplexer interface {

	// NewConn constructs a new connection
	// TODO rename to Wrap / Multiplex
	NewConn(c net.Conn, isServer bool) (MuxedConn, error)
}
