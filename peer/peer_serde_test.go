package peer_test

import (
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"
	. "github.com/libp2p/go-libp2p-core/test"
)

func TestPeerSerdePB(t *testing.T) {
	id, err := RandPeerID()
	if err != nil {
		t.Fatal(err)
	}
	b, err := id.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	var id2 peer.ID
	if err = id2.Unmarshal(b); err != nil {
		t.Fatal(err)
	}
	if id != id2 {
		t.Error("expected equal ids in circular serde test")
	}
}

func TestPeerSerdeJSON(t *testing.T) {
	id, err := RandPeerID()
	if err != nil {
		t.Fatal(err)
	}
	b, err := id.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	var id2 peer.ID
	if err = id2.UnmarshalJSON(b); err != nil {
		t.Fatal(err)
	}
	if id != id2 {
		t.Error("expected equal ids in circular serde test")
	}
}

func TestBinaryMarshaler(t *testing.T) {
	id, err := RandPeerID()
	if err != nil {
		t.Fatal(err)
	}
	b, err := id.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	var id2 peer.ID
	if err = id2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
	if id != id2 {
		t.Error("expected equal ids in circular serde test")
	}
}

func TestTextMarshaler(t *testing.T) {
	id, err := RandPeerID()
	if err != nil {
		t.Fatal(err)
	}
	b, err := id.MarshalText()
	if err != nil {
		t.Fatal(err)
	}

	var id2 peer.ID
	if err = id2.UnmarshalText(b); err != nil {
		t.Fatal(err)
	}
	if id != id2 {
		t.Error("expected equal ids in circular serde test")
	}
}
