package mux

import (
	"errors"
	"testing"
)

func TestUnwrapStreamError(t *testing.T) {
	foo := errors.New("foo")
	sce := &StreamCloseError{
		Code:   0,
		Reason: foo,
	}

	err := errors.Unwrap(sce)
	if err != foo {
		t.Fatalf("expected unwrapped error to be: %s; was: %s", foo, err)
	}
}
