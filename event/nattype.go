package event

import "github.com/libp2p/go-libp2p-core/network"

// NATDeviceProtocol is the transport protocol for which the NAT Device Type has been determined.
type NATDeviceProtocol int

const (
	// NATDeviceUDP means that the NAT Device Type has been determined for the UDP Protocol.
	NATDeviceUDP NATDeviceProtocol = iota
	// NATDeviceTCP means that the NAT Device Type has been determined for the TCP Protocol.
	NATDeviceTCP
)

// EvtNATDeviceTypeChanged is an event struct to be emitted when the type of the NAT device changes.
type EvtNATDeviceTypeChanged struct {
	Protocol      NATDeviceProtocol
	NatDeviceType network.NATDeviceType
}
