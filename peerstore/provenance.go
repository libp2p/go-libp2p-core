package peerstore

import (
	"fmt"
	"sync"
)

// AddressProvenance specifies the provenance of an address that has been stored
// in the peerstore.
//
// Values [0x00, 0x80) are RESERVED for the developers of libp2p.
//
// Users are free to register custom application-level provenances in the range
// [0x80, 0xff], via the RegisterAddressProvenance function.
type AddressProvenance uint8

const (
	// ProvenanceUnknown means that the provenance of an address is unknown.
	//
	// This is usually the case when an address has been inserted through the
	// legacy methods of the peerstore.
	ProvenanceUnknown = AddressProvenance(0x00)

	// ProvenanceThirdParty indicates that we've learnt this address via a third
	// party.
	ProvenanceThirdParty = AddressProvenance(0x10)

	// ProvenanceUntrusted means that an address has been returned by the peer
	// in question, but it is not authenticated, i.e. it is not part of a peer
	// record.
	ProvenanceUntrusted = AddressProvenance(0x20)

	// ProvenanceTrusted means that the address is part of an authenticated
	// standard libp2p peer record.
	ProvenanceTrusted = AddressProvenance(0x30)

	// ProvenanceManual means that an address has been specified manually by a
	// human, and therefore should take high precedence.
	ProvenanceManual = AddressProvenance(0x7f)
)

var (
	// provenanceLk only guards registration (i.e. writes, not reads).
	//
	// It is assumed that user-defined provenances will be registered during
	// system initialisation. Once the system becomes operational, the
	// descriptions are accessed by the String() method without taking a lock.
	provenanceLk sync.Mutex

	// provenanceDescs are the descriptions of the provenances, for debugging
	// purposes.
	provenanceDescs = [256]string{
		0x00: "unknown",
		0x10: "third_party",
		0x20: "untrusted",
		0x30: "trusted",
		0x7f: "manual",
	}
)

func (p AddressProvenance) String() string {
	return provenanceDescs[p]
}

func RegisterAddressProvenance(p AddressProvenance, desc string) error {
	provenanceLk.Lock()
	defer provenanceLk.Unlock()

	if p < 0x80 {
		return fmt.Errorf("failed to register user-defined address provenance "+
			"due to range violation; should be in [0x80,0xff]; was: %x", p)
	}
	if d := provenanceDescs[p]; d != "" {
		return fmt.Errorf("an address provenance for code %x already exists: %s", p, d)
	}
	provenanceDescs[p] = desc
	return nil
}
