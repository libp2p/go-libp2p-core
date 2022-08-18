package pnet

import "github.com/libp2p/go-libp2p/core/pnet"

// ErrNotInPrivateNetwork is an error that should be returned by libp2p when it
// tries to dial with ForcePrivateNetwork set and no PNet Protector
// Deprecated: use github.com/libp2p/go-libp2p/core/pnet.ErrNotInPrivateNetwork instead
var ErrNotInPrivateNetwork = pnet.ErrNotInPrivateNetwork

// Error is error type for ease of detecting PNet errors
// Deprecated: use github.com/libp2p/go-libp2p/core/pnet.Error instead
type Error = pnet.Error

// NewError creates new Error
// Deprecated: use github.com/libp2p/go-libp2p/core/pnet.NewError instead
func NewError(err string) error {
	return pnet.NewError(err)
}

// IsPNetError checks if given error is PNet Error
// Deprecated: use github.com/libp2p/go-libp2p/core/pnet.IsPNetError instead
func IsPNetError(err error) bool {
	return pnet.IsPNetError(err)
}
