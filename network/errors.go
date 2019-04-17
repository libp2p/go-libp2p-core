package network

import "errors"

// There are no addresses associated with a peer when they were needed.
var ErrNoRemoteAddrs = errors.New("no remote addresses")

// ErrNoConn is returned when attempting to open a stream to a peer with the NoDial
// option and no usable connection is available.
var ErrNoConn = errors.New("no usable connection to peer")
