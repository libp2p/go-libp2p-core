// Deprecated: This package has moved into go-libp2p as a sub-package: github.com/libp2p/go-libp2p/core/test.
package test

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/test"

	ma "github.com/multiformats/go-multiaddr"
)

// Deprecated: use github.com/libp2p/go-libp2p/core/test.GenerateTestAddrs instead
func GenerateTestAddrs(n int) []ma.Multiaddr {
	return test.GenerateTestAddrs(n)
}

// Deprecated: use github.com/libp2p/go-libp2p/core/test.AssertAddressesEqual instead
func AssertAddressesEqual(t *testing.T, exp, act []ma.Multiaddr) {
	test.AssertAddressesEqual(t, exp, act)
}
