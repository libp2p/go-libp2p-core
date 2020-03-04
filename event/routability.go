package event

import (
	"github.com/libp2p/go-libp2p-core/network"
)

// EvtLocalRoutability is an event struct to be emitted with the local's node
// routability changes state.
//
// This event is usually emitted by the AutoNAT subsystem.
type EvtLocalRoutability struct {
	Routability network.Routability
}
