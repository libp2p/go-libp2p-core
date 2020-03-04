package peerstore

import (
	"fmt"
	"sync"
)

// Label is an arbitrary label attached to an address or a peer. A Label can
// indicate the source of an address, a "pin" by a subsystem, an address scope,
// or an arbitrary tag.
//
// Labels are byte values in binary format, bound to strings via the `BindLabel`
// function to facilitate debugging.
//
// Values [0x00, 0x80) are RESERVED for the developers of libp2p. Users are free
// to register custom application-level labels in the range [0x80, 0xff], via
// the BindLabel function.
type Label uint8

const (
	// LabelSourceUnknown means that the source of an address is unknown.
	//
	// This is usually the case when an address has been inserted through the
	// legacy methods of the peerstore.
	LabelSourceUnknown = Label(0x00)

	// LabelSourceThirdParty indicates that we've learnt this address via a
	// third party.
	LabelSourceThirdParty = Label(0x01)

	// LabelSourceUncertified means that an address has been returned by the
	// peer in question, but it is not authenticated, i.e. it is not part of a
	// certified peer record.
	LabelSourceUncertified = Label(0x02)

	// LabelSourceCertified means that the address is part of a standard
	// certified libp2p peer record.
	LabelSourceCertified = Label(0x03)

	// LabelSourceManual means that an address has been specified manually by a
	// human, and therefore should take high precedence.
	LabelSourceManual = Label(0x04)
)

var (
	// labelsLk only guards registration (i.e. writes, not reads).
	//
	// It is assumed that user-defined labels will be registered during system
	// initialisation. Once the system becomes operational, the descriptions are
	// accessed by the String() method *without taking a lock*.
	labelsLk sync.Mutex

	// labelsTable stores label allocations and their string mappings.
	labelsTable = [256]string{
		0x00: "unknown",
		0x01: "third_party",
		0x02: "untrusted",
		0x03: "trusted",
		0x04: "manual",
	}
)

func (p Label) String() string {
	return labelsTable[p]
}

// BindLabel registers a user-defined label and binds it to a string
// description.
func BindLabel(p Label, desc string) error {
	labelsLk.Lock()
	defer labelsLk.Unlock()

	if p < 0x80 {
		return fmt.Errorf("failed to register user-defined peerstore label "+
			"due to range violation; should be in [0x80,0xff]; was: %x", p)
	}
	if d := labelsTable[p]; d != "" {
		return fmt.Errorf("a label mapping for %x already exists: %s", p, d)
	}
	labelsTable[p] = desc
	return nil
}
