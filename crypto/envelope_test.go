package crypto_test

import (
	"bytes"
	"github.com/gogo/protobuf/proto"
	. "github.com/libp2p/go-libp2p-core/crypto"
	pb "github.com/libp2p/go-libp2p-core/crypto/pb"
	"github.com/libp2p/go-libp2p-core/test"
	"testing"
)

// Make an envelope, verify & open it, marshal & unmarshal it
func TestEnvelopeHappyPath(t *testing.T) {
	priv, pub, err := test.RandTestKeyPair(Ed25519, 256)
	if err != nil {
		t.Error(err)
	}
	contents := []byte("happy hacking")
	domain := "libp2p-testing"
	typeHint := []byte("/libp2p/testdata")
	envelope, err := MakeEnvelope(priv, domain, typeHint, contents)
	if err != nil {
		t.Errorf("error constructing envelope: %v", err)
	}

	if !envelope.PublicKey.Equals(pub) {
		t.Error("envelope has unexpected public key")
	}

	if bytes.Compare(typeHint, envelope.TypeHint) != 0 {
		t.Error("TypeHint does not match typeHint used to construct envelope")
	}

	valid, err := envelope.Validate(domain)
	if err != nil {
		t.Errorf("error validating envelope: %v", err)
	}
	if !valid {
		t.Error("envelope should be valid, but Valid returns false")
	}

	c, err := envelope.Open(domain)
	if err != nil {
		t.Errorf("error opening envelope: %v", err)
	}
	if bytes.Compare(c, contents) != 0 {
		t.Error("contents of envelope do not match input")
	}

	serialized, err := envelope.Marshal()
	if err != nil {
		t.Errorf("error serializing envelope: %v", err)
	}
	deserialized, err := UnmarshalEnvelope(serialized)
	if err != nil {
		t.Errorf("error deserializing envelope: %v", err)
	}

	if !envelope.Equals(deserialized) {
		t.Error("round-trip serde results in unequal envelope structures")
	}
}

func TestEnvelopeValidateFailsForDifferentDomain(t *testing.T) {
	priv, _, err := test.RandTestKeyPair(Ed25519, 256)
	if err != nil {
		t.Error(err)
	}
	contents := []byte("happy hacking")
	domain := "libp2p-testing"
	typeHint := []byte("/libp2p/testdata")
	envelope, err := MakeEnvelope(priv, domain, typeHint, contents)
	if err != nil {
		t.Errorf("error constructing envelope: %v", err)
	}

	valid, err := envelope.Validate("other-domain")
	if err != nil {
		t.Errorf("error validating envelope: %v", err)
	}
	if valid {
		t.Error("envelope should be invalid, but Valid returns true")
	}
}

func TestEnvelopeValidateFailsIfTypeHintIsAltered(t *testing.T) {
	priv, _, err := test.RandTestKeyPair(Ed25519, 256)
	if err != nil {
		t.Error(err)
	}
	contents := []byte("happy hacking")
	domain := "libp2p-testing"
	typeHint := []byte("/libp2p/testdata")
	envelope, err := MakeEnvelope(priv, domain, typeHint, contents)
	if err != nil {
		t.Errorf("error constructing envelope: %v", err)
	}
	envelope.TypeHint = []byte("foo")
	valid, err := envelope.Validate("other-domain")
	if err != nil {
		t.Errorf("error validating envelope: %v", err)
	}
	if valid {
		t.Error("envelope should be invalid, but Valid returns true")
	}
}

func TestEnvelopeValidateFailsIfContentsAreAltered(t *testing.T) {
	priv, _, err := test.RandTestKeyPair(Ed25519, 256)
	if err != nil {
		t.Error(err)
	}
	contents := []byte("happy hacking")
	domain := "libp2p-testing"
	typeHint := []byte("/libp2p/testdata")
	envelope, err := MakeEnvelope(priv, domain, typeHint, contents)
	if err != nil {
		t.Errorf("error constructing envelope: %v", err)
	}

	// since the contents field is private, we'll serialize and alter the
	// serialized protobuf
	serialized, err := envelope.Marshal()
	if err != nil {
		t.Errorf("error serializing envelope: %v", err)
	}

	msg := pb.SignedEnvelope{}
	err = proto.Unmarshal(serialized, &msg)
	if err != nil {
		t.Errorf("error deserializing envelope: %v", err)
	}
	msg.Contents = []byte("totally legit, trust me")
	serialized, err = proto.Marshal(&msg)

	// unmarshal our altered envelope
	deserialized, err := UnmarshalEnvelope(serialized)
	if err != nil {
		t.Errorf("error deserializing envelope: %v", err)
	}

	// verify that it's now invalid
	valid, err := deserialized.Validate(domain)
	if err != nil {
		t.Errorf("error validating envelope: %v", err)
	}
	if valid {
		t.Error("envelope should be invalid, but Valid returns true")
	}
}
