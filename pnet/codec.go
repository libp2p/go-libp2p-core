package pnet

import (
	"io"

	"github.com/libp2p/go-libp2p/core/pnet"
)

// DecodeV1PSK reads a Multicodec encoded V1 PSK.
// Deprecated: use github.com/libp2p/go-libp2p/core/pnet.DecodeV1PSK instead
func DecodeV1PSK(in io.Reader) (PSK, error) {
	return pnet.DecodeV1PSK(in)
}
