package event

import "github.com/libp2p/go-libp2p-core/network"

// EvtNATDeviceTypeChanged is an event struct to be emitted when the type of the NAT device changes.
type EvtNATDeviceTypeChanged struct {
	NatDeviceType network.NATDeviceType
}
