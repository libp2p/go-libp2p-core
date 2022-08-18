package pnet

import (
	"github.com/libp2p/go-libp2p/core/pnet"
)

// EnvKey defines environment variable name for forcing usage of PNet in libp2p
// When environment variable of this name is set to "1" the ForcePrivateNetwork
// variable will be set to true.
// Deprecated: use github.com/libp2p/go-libp2p/core/pnet.EnvKey instead
const EnvKey = pnet.EnvKey

// ForcePrivateNetwork is boolean variable that forces usage of PNet in libp2p
// Setting this variable to true or setting LIBP2P_FORCE_PNET environment variable
// to true will make libp2p to require private network protector.
// If no network protector is provided and this variable is set to true libp2p will
// refuse to connect.
// Deprecated: use github.com/libp2p/go-libp2p/core/pnet.ForcePrivateNetwork instead
var ForcePrivateNetwork = pnet.ForcePrivateNetwork
