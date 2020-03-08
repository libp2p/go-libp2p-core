// Package backoff provides facilities for working with interdependent backoff timers.
//
// This package is motivated by https://github.com/libp2p/go-libp2p-swarm/issues/37.
//
// We would like to associate backoff logic to each component of a multiaddr.
// For instance, consider a typical multiaddr of the form IP/transport/protocol.
// The backoff associated with IP dictates whether any multiaddr
// using that IP can be dialed (regardless of transport and protocol).
// The backoff associated with IP and transport dictates whether
// any multiaddr using that IP and transport can be dialed (regardless of protocol).
// And so on.
//
// Furthermore, we stipulate that software components that utilize backoffs on multiaddrs
// should share the backoff timers in order to benefit from each other's learnings.
//
// To this end, this package provides the definition of a backoff timer (as the interface BackoffTimer),
// as well as an interface for a shared set of backoff timers, called SharedBackoffs.
package backoff

import (
	"time"

	ma "github.com/multiformats/go-multiaddr"
)

// SharedBackoffs maintains backoff timers that are shared by multiple IPFS components.
// Some of the timers are correlated. If this is the case, the relationships are documented in the access methods below.
type SharedBackoffs interface {

	// Timers for dialing multiaddresses.

	// IP returns the timer associated with a given IP address.
	// If the argument multiaddr does not start with an IP component, a panic is thrown.
	// The IP timer should be backed off, if it is determined that the IP itself is unreachable
	// (e.g. in response to a "no route to IP" error).
	// The IP timer should be cleared whenever any successful connection to this IP is established.
	IP(ma.Multiaddr) BackoffTimer

	// IPTransport returns the timer associated with a given IP address and a transport (TCP, UDP, etc.).
	// If the argument multiaddr does not start with an IP and transport components, a panic is thrown.
	// The IP/Transport timer should be backed off, if it is determined that the IP is reachable but the transport is not.
	// The IP/Transport timer should be cleared whenever any successful connection to this IP/transport pair is established.
	// The IP/Transport timer will report a clear state whenever the IP and  IP/Transport timers are both clear.
	IPTransport(ma.Multiaddr) BackoffTimer

	// IPTransportSwarm returns the timer associated with a given IP address, a transport (TCP, UDP, etc.) and the swarm service.
	// If the argument multiaddr does not start with an IP and transport components, a panic is thrown.
	// The IP/Transport/Swarm timer should be backed off, if it is determined that the IP is reachable
	// using the given transport, but the swarm service is not supported.
	// The IP/Transport/Swarm timer should be cleared whenever any successful connection to the IP/transport pair is established
	// and the swarm service is available.
	// The IP/Transport/Swarm timer will report a clear state whenever the respective IP, IP/Transport and IP/Transport/Swarm
	// timers are all clear.
	IPTransportSwarm(ma.Multiaddr) BackoffTimer

	// IPTransportSwarmProtocol returns the timer associated with a given IP address, transport (TCP, UDP, etc.) and
	// protocol, supported through the swarm service.
	// If the argument multiaddr does not start with an IP, transport and protocol components, a panic is thrown.
	// The IP/Transport/Swarm/Protocol timer should be backed off, if it is determined that the IP is reachable
	// using the given transport, the swarm service is supported, but the protocol is not.
	// The IP/Transport/Swarm/Protocol timer should be cleared whenever any successful connection to the protocol
	// is established through the swarm service on the given IP/transport.
	// The IP/Transport/Swarm/Protocol timer will report a clear state whenever the respective
	// IP, IP/Transport and IP/Transport/Swarm, and IP/Transport/Swarm/Protocol timers are all clear.
	IPTransportSwarmProtocol(ma.Multiaddr) BackoffTimer

	// Other shared backoff timers should be added here.
}

// BackoffTimer is a user-facing interface to a backoff timer.
type BackoffTimer interface {
	// Wait blocks until the timer and all of its ancestors (parent, grandparent, etc.), if any, have been cleared.
	Wait()
	// TimeToClear returns the duration remaining until the back off state is cleared.
	// Zero or negative durations indicate that the state is already cleared.
	TimeToClear(now time.Time) time.Duration
	// Clear clears this timer and returns instantaneously.
	// For example, a user might call this function after a successul connection attempt.
	Clear()
	// Backoff places the timer in backing off state and returns instantaneously.
	// For example, a user might call this function after an unsuccessful connection attempt.
	Backoff()
}
