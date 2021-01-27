package network

// NATDeviceType indicates the type of the NAT device i.e. whether it is a Hard or an Easy NAT.
type NATDeviceType int

const (
	// NATDeviceTypeUnknown indicates that the type of the NAT device is unknown.
	NATDeviceTypeUnknown NATDeviceType = iota

	// NATDeviceTypeEasy indicates that the NAT device is an Easy NAT i.e. it supports consistent endpoint translation.
	// NAT traversal via hole punching is possible with this NAT type if the remote peer is also behind an Easy NAT.
	NATDeviceTypeEasy

	// NATDeviceTypeHard indicates that the NAT device is a Hard NAT that does NOT support
	// consistent endpoint translation.
	// NAT traversal via hole-punching is NOT possible with this NAT type irrespective of the remote peer's NAT type.
	NATDeviceTypeHard
)
