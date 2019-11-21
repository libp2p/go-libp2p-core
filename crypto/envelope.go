package crypto

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/gogo/protobuf/proto"

	"github.com/libp2p/go-buffer-pool"
	pb "github.com/libp2p/go-libp2p-core/crypto/pb"
)

// SignedEnvelope contains an arbitrary []byte payload, signed by a libp2p peer.
// Envelopes are signed in the context of a particular "domain", which is a string
// specified when creating and verifying the envelope. You must know the domain
// string used to produce the envelope in order to verify the signature and
// access the payload.
type SignedEnvelope struct {
	// The public key that can be used to verify the signature and derive the peer id of the signer.
	PublicKey PubKey

	// A binary identifier that indicates what kind of data is contained in the payload.
	// TODO(yusef): enforce multicodec prefix
	PayloadType []byte

	// The envelope payload.
	Payload []byte

	// The signature of the domain string, type hint, and payload.
	signature []byte
}

var errEmptyDomain = errors.New("envelope domain must not be empty")
var errInvalidSignature = errors.New("invalid signature or incorrect domain")

// MakeEnvelope constructs a new SignedEnvelope using the given privateKey.
//
// The required 'domain' string contextualizes the envelope for a particular purpose,
// and must be supplied when verifying the signature.
//
// The 'PayloadType' field indicates what kind of data is contained and may be empty.
func MakeEnvelope(privateKey PrivKey, domain string, payloadType []byte, payload []byte) (*SignedEnvelope, error) {
	if len(domain) == 0 {
		return nil, errEmptyDomain
	}
	toSign, err := makeSigBuffer(domain, payloadType, payload)
	if err != nil {
		return nil, err
	}
	sig, err := privateKey.Sign(toSign)
	if err != nil {
		return nil, err
	}

	return &SignedEnvelope{
		PublicKey:   privateKey.GetPublic(),
		PayloadType: payloadType,
		Payload:     payload,
		signature:   sig,
	}, nil
}

// OpenEnvelope unmarshals a serialized SignedEnvelope, validating its signature
// using the provided 'domain' string.
func OpenEnvelope(envelopeBytes []byte, domain string) (*SignedEnvelope, error) {
	e, err := UnmarshalEnvelopeWithoutValidating(envelopeBytes)
	if err != nil {
		return nil, err
	}
	err = e.validate(domain)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// UnmarshalEnvelopeWithoutValidating unmarshals a serialized SignedEnvelope protobuf message,
// without validating its contents. Should not be used unless you have a very good reason
// (e.g. testing)!
func UnmarshalEnvelopeWithoutValidating(serializedEnvelope []byte) (*SignedEnvelope, error) {
	var e pb.SignedEnvelope
	if err := proto.Unmarshal(serializedEnvelope, &e); err != nil {
		return nil, err
	}
	key, err := PublicKeyFromProto(e.PublicKey)
	if err != nil {
		return nil, err
	}
	return &SignedEnvelope{
		PublicKey:   key,
		PayloadType: e.PayloadType,
		Payload:     e.Payload,
		signature:   e.Signature,
	}, nil
}

// Marshal returns a byte slice containing a serailized protobuf representation of a SignedEnvelope.
func (e *SignedEnvelope) Marshal() ([]byte, error) {
	key, err := PublicKeyToProto(e.PublicKey)
	if err != nil {
		return nil, err
	}
	msg := pb.SignedEnvelope{
		PublicKey:   key,
		PayloadType: e.PayloadType,
		Payload:     e.Payload,
		Signature:   e.signature,
	}
	return proto.Marshal(&msg)
}

// Equal returns true if the other SignedEnvelope has the same
// public key, payload, payload type, and signature. This
// implies that they were also created with the same domain string.
func (e *SignedEnvelope) Equal(other *SignedEnvelope) bool {
	return e.PublicKey.Equals(other.PublicKey) &&
		bytes.Compare(e.PayloadType, other.PayloadType) == 0 &&
		bytes.Compare(e.Payload, other.Payload) == 0 &&
		bytes.Compare(e.signature, other.signature) == 0
}

// validate returns true if the envelope signature is valid for the given 'domain',
// or false if it is invalid. May return an error if signature validation fails.
func (e *SignedEnvelope) validate(domain string) error {
	toVerify, err := makeSigBuffer(domain, e.PayloadType, e.Payload)
	if err != nil {
		return err
	}
	valid, err := e.PublicKey.Verify(toVerify, e.signature)
	if err != nil {
		return err
	}
	if !valid {
		return errInvalidSignature
	}
	return nil
}

// makeSigBuffer is a helper function that prepares a buffer to sign or verify.
func makeSigBuffer(domain string, payloadType []byte, payload []byte) ([]byte, error) {
	domainBytes := []byte(domain)
	fields := [][]byte{domainBytes, payloadType, payload}

	const lengthPrefixSize = 8
	size := 0
	for _, f := range fields {
		size += len(f) + lengthPrefixSize
	}

	b := pool.NewBuffer(nil)
	b.Grow(size)

	for _, f := range fields {
		err := writeField(b, f)
		if err != nil {
			return nil, err
		}
	}

	return b.Bytes(), nil
}

func writeField(b *pool.Buffer, f []byte) error {
	_, err := b.Write(encodedSize(f))
	if err != nil {
		return err
	}
	_, err = b.Write(f)
	return err
}

func encodedSize(content []byte) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(len(content)))
	return b
}
