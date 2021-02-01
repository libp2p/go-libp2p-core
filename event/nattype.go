package event

import "github.com/libp2p/go-libp2p-core/network"

// NATTransportProtocol is the transport protocol for which the NAT Device Type has been determined.
type NATTransportProtocol int

const (
	// NATTransportUDP means that the NAT Device Type has been determined for the UDP Protocol.
	NATTransportUDP NATTransportProtocol = iota
	// NATTransportTCP means that the NAT Device Type has been determined for the TCP Protocol.
	NATTransportTCP
)

func (n NATTransportProtocol) String() string {
	switch n {
	case 0:
		return "UDP"
	case 1:
		return "TCP"
	default:
		return "unrecognized"
	}
}

// EvtNATDeviceTypeChanged is an event struct to be emitted when the type of the NAT device changes for a Transport Protocol.
//
// Note: This event is meaningful ONLY if the AutoNAT Reachability is Private.
// Consumers of this event should ALSO consume the `EvtLocalReachabilityChanged` event and interpret
// this event ONLY if the Reachability on the `EvtLocalReachabilityChanged` is Private.
type EvtNATDeviceTypeChanged struct {
	// TransportProtocol is the Transport Protocol for which the NAT Device Type has been determined.
	TransportProtocol NATTransportProtocol
	// NatDeviceType indicates the type of the NAT Device for the Transport Protocol.
	// Currently, it can be either a `Cone NAT` or a `Symmetric NAT`. Please see the detailed documentation
	// on `network.NATDeviceType` enumerations for a better understanding of what these types mean and
	// how they impact Connectivity and Hole Punching.
	NatDeviceType network.NATDeviceType
}
