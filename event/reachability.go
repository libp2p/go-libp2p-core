package event

import (
	"github.com/libp2p/go-libp2p/core/event"
)

// EvtLocalReachabilityChanged is an event struct to be emitted when the local's
// node reachability changes state.
//
// This event is usually emitted by the AutoNAT subsystem.
// Deprecated: use github.com/libp2p/go-libp2p/core/event.EvtLocalReachabilityChanged instead
type EvtLocalReachabilityChanged = event.EvtLocalReachabilityChanged
