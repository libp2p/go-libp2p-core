package mux

import (
	"fmt"
	"sync/atomic"
)

// StreamCloseError encapsulates error information that the application may
// want to send to the peer when closing a stream.
//
// This data will be used on a best-effort basis. Not all stream muxers support
// rich error signalling semantics.
//
// To check if the underlying muxer has effectively used the error details,
// call the Used() method.
//
// Authors/maintainers of stream muxers can signal that the error has been used
// by calling MarkUsed().
type StreamCloseError struct {
	Code   int
	Reason error

	// for alignment, this field is last, as the above two fields are
	// likely one word each.
	used int32
}

// Used returns whether the stream muxer has used any of the details of this
// error.
func (e *StreamCloseError) Used() bool {
	return atomic.LoadInt32(&e.used) == 1
}

// MarkUsed allows a stream muxer to mark the details of this error as used.
func (e *StreamCloseError) MarkUsed() {
	atomic.StoreInt32(&e.used, 1)
}

func (e StreamCloseError) Unwrap() error {
	return e.Reason
}

func (e StreamCloseError) Error() string {
	msg := fmt.Sprintf("stream close error; code: %d", e.Code)
	if e.Reason != nil {
		msg = msg + "; err: " + e.Reason.Error()
	}
	return msg
}
