package test

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/test"
)

// Deprecated: use github.com/libp2p/go-libp2p/core/test.AssertNilError instead
func AssertNilError(t *testing.T, err error) {
	test.AssertNilError(t, err)
}

// Deprecated: use github.com/libp2p/go-libp2p/core/test.ExpectError instead
func ExpectError(t *testing.T, err error, msg string) {
	test.ExpectError(t, err, msg)
}
