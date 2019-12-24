package connmgr

import (
	"io"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
)

// Decayer is implemented by connection managers supporting decaying tags. A
// decaying tag is one whose value automatically decays over time.
//
// The actual application of the decay behaviour is encapsulated in a
// user-provided decaying function (DecayFn). The function is called on every
// tick (determined by the interval parameter), and returns either the new value
// of the tag, or whether it should be erased altogether.
//
// We do not set values on a decaying tag. Rather, we "bump" decaying tags by a
// delta. This calls the BumpFn with the old value and the delta, to determine
// the new value.
//
// Such a pluggable design affords a great deal of flexibility and versatility.
// Behaviours that are straightfoward to implement include:
//
//  * Decay a tag by -1, or by half its current value, on every tick.
//  * Every time a value is bumped, sum it to its current value.
//  * Exponentially boost a score with every bump.
//  * Sum the incoming score, but keep it within min, max bounds.
//
// Commonly used DecayFns and BumpFns are provided in the go-libp2p-connmgr
// module.
type Decayer interface {
	io.Closer

	// RegisterDecayingTag creates and registers a new decaying tag, if and only
	// if a tag with the supplied name doesn't exist yet. Otherwise, an error is
	// returned.
	//
	// The caller provides the interval at which the tag is refreshed, as well
	// as the decay function and the bump function. Refer to godocs on DecayFn
	// and BumpFn for more info.
	RegisterDecayingTag(name string, interval time.Duration, decayFn DecayFn, bumpFn BumpFn) (DecayingTag, error)
}

// DecayFn applies a decay to the peer's score. The implementation must call
// DecayFn at the interval supplied when registering the tag.
//
// It receives a copy of the decaying value, and returns the score after
// applying the decay, as well as a flag to signal if the tag should be erased.
type DecayFn func(value DecayingValue) (after int, rm bool)

// BumpFn applies a delta onto an existing score, and returns the new score.
//
// Non-trivial bump functions include exponential boosting, moving averages,
// ceilings, etc.
type BumpFn func(value DecayingValue, delta int) (after int)

// DecayingTag represents a decaying tag. The tag is a long-lived general
// object, used to operate on tag values for peers.
type DecayingTag interface {

	// Bump applies a delta to a tag value, calling its bump function. The bump
	// may be applied asynchronously, in which case the returned error is used
	// to indicate an anomaly when queuing.
	Bump(peer peer.ID, delta int) error
}

// DecayingValue represents a value for a decaying tag.
type DecayingValue struct {
	// Tag points to the tag this value belongs to.
	Tag DecayingTag

	// Peer is the peer ID to whom this value is associated.
	Peer peer.ID

	// Added is the timestamp when this value was added for the first time for
	// a tag and a peer.
	Added time.Time

	// LastVisit is the timestamp of the last visit.
	LastVisit time.Time

	// Value is the current value of the tag.
	Value int
}
