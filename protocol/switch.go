// Package protocol provides core interfaces for protocol routing and negotiation in libp2p.
package protocol

import (
	"io"
)

// HandlerFunc is a user-provided function used by the Router to
// handle a protocol/stream.
//
// Will be invoked with the protocol ID string as the first argument,
// which may differ from the ID used for registration if the handler
// was registered using a match function.
type HandlerFunc = func(protocol string, rwc io.ReadWriteCloser) error

// Router is an interface that allows users to add and remove protocol handlers,
// which will be invoked when incoming stream requests for registered protocols
// are accepted.
type Router interface {

	// AddHandler registers the given handler to be invoked for
	// an exact literal match of the given protocol ID string.
	AddHandler(protocol string, handler HandlerFunc)

	// AddHandlerWithFunc registers the given handler to be invoked
	// for exact literal matches of the given protocol ID string,
	// **or** when the provided match function returns true.
	//
	// The match function will be invoked with an incoming protocol
	// ID string when the router is unable to find an exact literal
	// match among the registered handlers.
	AddHandlerWithFunc(protocol string, match func(string) bool, handler HandlerFunc)

	// RemoveHandler removes the registered handler (if any) for the
	// given protocol ID string.
	RemoveHandler(protocol string)

	// Protocols returns a list of all registered protocol ID strings.
	// Note that the Router may be able to handle protocol IDs not
	// included in this list if handlers were added with match functions
	// using AddHandlerWithFunc.
	Protocols() []string
}

// Negotiator is a component capable of reaching agreement over what protocols
// to use for a given stream of communication.
//
// The Negotiator is responsible for proposing the protocol to use for a given
// stream, and it responds to proposals sent by the other end of the stream.
type Negotiator interface {
	NegotiateLazy(rwc io.ReadWriteCloser) (io.ReadWriteCloser, string, HandlerFunc, error)
	Negotiate(rwc io.ReadWriteCloser) (string, HandlerFunc, error)
	Handle(rwc io.ReadWriteCloser) error
}

// Switch is the component responsible for "dispatching" incoming stream requests to
// their corresponding stream handlers. It is both a Negotiator and a Router.
type Switch interface {
	Router
	Negotiator
}
