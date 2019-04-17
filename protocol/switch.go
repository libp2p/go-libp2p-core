// Package protocol provides core interfaces for protocol routing and negotiation in libp2p.
package protocol

import (
	"io"
)

// HandlerFunc is a user-provided function used by the Router to
// handle a protocol/stream.
type HandlerFunc = func(protocol string, rwc io.ReadWriteCloser) error

type Router interface {
	AddHandler(protocol string, handler HandlerFunc)
	AddHandlerWithFunc(protocol string, match func(string) bool, handler HandlerFunc)
	RemoveHandler(protocol string)
	Protocols() []string
}

type Negotiator interface {
	NegotiateLazy(rwc io.ReadWriteCloser) (io.ReadWriteCloser, string, HandlerFunc, error)
	Negotiate(rwc io.ReadWriteCloser) (string, HandlerFunc, error)
	Handle(rwc io.ReadWriteCloser) error
}

type Switch interface {
	Router
	Negotiator
}
