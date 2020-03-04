// Package backoff provides facilities for creating interdependent backoff timers.
//
// The design of "backoff trees" proposed in this package is motivated by
// https://github.com/libp2p/go-libp2p-swarm/issues/37.
//
// The design proposed here can be dubbed "backoff trees", and is motivated by the following use case in libp2p:
//
// (1) Each peer has an associated backoff timer, whose significance is that no interaction
// with this peer can commence unless the timer is clear.
// (2) Each peer has a subordinate "find peer" timer, such that find peer operations should not
// be attempted unless the "find peer" timer is clear and the "peer" timer is clear, as well.
// (3) Each peer has any number of subordinate "transport" timers, such that dial attempts on any
// given transport should not happen unless the "transport" timer is clear and the "peer" timer is clear, as well.
//	(4) Each TCP transport may have subordinate "IP" timers, such that TCP dial attempts to the given IP
// should not commence unless the "IP", "transport" and "peer" timers are all clear.
// And so on.
//
// This example highlights a hierarchical relationship that one might express, for instance as, as:
//		peer operations
//			find operations
//			transport-related operations
//				TCP-related operations
//					IP-related operations
//						Port-related operations
//				QUIC-related operations
//
// All nodes in this hierarchy may wish to use a different backoff policy.
// However, in all cases, a timer node is considered "clear" if its own backoff timer is clear,
// as well as the timers of all of its ancestors (all the way to the root).
//
// REMARKS
//
// Note that this design can be extended to "backoff DAGS", wherein a timer can have multiple parents.
// This can be useful, for instance, in the case when an IP address is viewed as a subordinate timer to multiple protocols,
// like TCP and QUIC.
//
// While supporting timers with multiple parents is straightforward, it is not clear that it can be used conveniently
// by independent code paths. In particular, a runtime instance of the QUIC transport may not (and probably should not)
// be aware that a TCP transport instance is running in parallel.
package backoff

import (
	"time"
)

// NewBackoffTimer creates a new backoff timer.
// The timer name is for display purposes only.
NewBackoffTimer(name string, policy BackoffPolicy, parent BackoffTimer) BackoffTimer

// BackoffPolicy represents a specific backoff logic.
//
// Implementations of BackoffPolicy are purely concerned with the "arithmetic"
// of computing when the respective timer should be cleared (e.g. for making new connection retries).
//
// This interface allows for the implementation of flexible backoff policies.
// For instance, a policy could treat a burst of backoffs as a single one.
type BackoffPolicy interface {

	// Clear informs the policy of the current time and sets its state to cleared.
	Clear(now time.Time)

	// Backoff informs the policy of the current time and sets its state to backing off.
	Backoff(now time.Time)

	// TimeToClear informs the policy of the current time and returns the duration
	// remaining until the back off state is cleared. Zero or negative durations indicate
	// that the state is already cleared.
	TimeToClear(now time.Time) time.Duration
}

// BackoffTimer is a synchronous user-facing interface to a backoff timer.
// Timers are created using NewBackoffTimer.
type BackoffTimer interface {
	// Wait blocks until the timer and all of its ancestors (parent, grandparent, etc.), if any, have been cleared.
	Wait()
	// Clear clears this timer and returns instantaneously.
	// For example, a user might call this function after a successul connection attempt.
	Clear()
	// Backoff places the timer in backing off state and returns instantaneously.
	// For example, a user might call this function after an unsuccessful connection attempt.
	Backoff()
}
