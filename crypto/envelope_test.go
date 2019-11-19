package crypto_test

import (
	"bytes"
	"github.com/golang/protobuf/proto"

	. "github.com/libp2p/go-libp2p-core/crypto"
	pb "github.com/libp2p/go-libp2p-core/crypto/pb"
	"github.com/libp2p/go-libp2p-core/test"
	"testing"
)

// Make an envelope, verify & open it, marshal & unmarshal it
func TestEnvelopeHappyPath(t *testing.T) {
	priv, pub, err := test.RandTestKeyPair(Ed25519, 256)
	test.AssertNilError(t, err)
	payload := []byte("happy hacking")
	domain := "libp2p-testing"
	payloadType := []byte("/libp2p/testdata")
	envelope, err := MakeEnvelope(priv, domain, payloadType, payload)
	test.AssertNilError(t, err)

	if !envelope.PublicKey().Equals(pub) {
		t.Error("envelope has unexpected public key")
	}

	if bytes.Compare(payloadType, envelope.PayloadType()) != 0 {
		t.Error("PayloadType does not match payloadType used to construct envelope")
	}

	serialized, err := envelope.Marshal()
	test.AssertNilError(t, err)
	deserialized, err := OpenEnvelope(serialized, domain)
	test.AssertNilError(t, err)

	if bytes.Compare(deserialized.Payload(), payload) != 0 {
		t.Error("payload of envelope does not match input")
	}

	if !envelope.Equal(deserialized) {
		t.Error("round-trip serde results in unequal envelope structures")
	}
}

func TestMakeEnvelopeFailsWithEmptyDomain(t *testing.T) {
	priv, _, err := test.RandTestKeyPair(Ed25519, 256)
	if err != nil {
		t.Error(err)
	}
	payload := []byte("happy hacking")
	payloadType := []byte("/libp2p/testdata")
	_, err = MakeEnvelope(priv, "", payloadType, payload)
	test.ExpectError(t, err, "making an envelope with an empty domain should fail")
}

func TestEnvelopeValidateFailsForDifferentDomain(t *testing.T) {
	priv, _, err := test.RandTestKeyPair(Ed25519, 256)
	test.AssertNilError(t, err)
	payload := []byte("happy hacking")
	domain := "libp2p-testing"
	payloadType := []byte("/libp2p/testdata")
	envelope, err := MakeEnvelope(priv, domain, payloadType, payload)
	test.AssertNilError(t, err)
	serialized, err := envelope.Marshal()
	// try to open our modified envelope
	_, err = OpenEnvelope(serialized, "wrong-domain")
	test.ExpectError(t, err, "should not be able to open envelope with incorrect domain")
}

func TestEnvelopeValidateFailsIfTypeHintIsAltered(t *testing.T) {
	priv, _, err := test.RandTestKeyPair(Ed25519, 256)
	test.AssertNilError(t, err)
	payload := []byte("happy hacking")
	domain := "libp2p-testing"
	payloadType := []byte("/libp2p/testdata")
	envelope, err := MakeEnvelope(priv, domain, payloadType, payload)
	test.AssertNilError(t, err)
	serialized := alterMessageAndMarshal(t, envelope, func(msg *pb.SignedEnvelope) {
		msg.PayloadType = []byte("foo")
	})
	// try to open our modified envelope
	_, err = OpenEnvelope(serialized, domain)
	test.ExpectError(t, err, "should not be able to open envelope with modified payloadType")
}

func TestEnvelopeValidateFailsIfContentsAreAltered(t *testing.T) {
	priv, _, err := test.RandTestKeyPair(Ed25519, 256)
	test.AssertNilError(t, err)
	payload := []byte("happy hacking")
	domain := "libp2p-testing"
	payloadType := []byte("/libp2p/testdata")
	envelope, err := MakeEnvelope(priv, domain, payloadType, payload)
	test.AssertNilError(t, err)

	serialized := alterMessageAndMarshal(t, envelope, func(msg *pb.SignedEnvelope) {
		msg.Payload = []byte("totally legit, trust me")
	})
	// try to open our modified envelope
	_, err = OpenEnvelope(serialized, domain)
	test.ExpectError(t, err, "should not be able to open envelope with modified payload")
}

// Since we're outside of the crypto package (to avoid import cycles with test package),
// we can't alter the fields in a SignedEnvelope directly. This helper marshals
// the envelope to a protobuf and calls the alterMsg function, which should
// alter the protobuf message.
// Returns the serialized altered protobuf message.
func alterMessageAndMarshal(t *testing.T, envelope *SignedEnvelope, alterMsg func(*pb.SignedEnvelope)) []byte {
	t.Helper()
	serialized, err := envelope.Marshal()
	test.AssertNilError(t, err)
	msg := pb.SignedEnvelope{}
	err = proto.Unmarshal(serialized, &msg)
	test.AssertNilError(t, err)
	alterMsg(&msg)
	serialized, err = msg.Marshal()
	test.AssertNilError(t, err)
	return serialized
}
