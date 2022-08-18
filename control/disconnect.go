// Deprecated: This package has moved into go-libp2p as a sub-package: github.com/libp2p/go-libp2p/core/control.
package control

import "github.com/libp2p/go-libp2p/core/control"

// DisconnectReason communicates the reason why a connection is being closed.
//
// A zero value stands for "no reason" / NA.
//
// This is an EXPERIMENTAL type. It will change in the future. Refer to the
// connmgr.ConnectionGater godoc for more info.
// Deprecated: use github.com/libp2p/go-libp2p/core/control instead
type DisconnectReason = control.DisconnectReason
