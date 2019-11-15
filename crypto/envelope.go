package crypto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/golang/protobuf/proto"
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

	// The envelope payload. This is private to discourage accessing the payload without verifying the signature.
	// To access, use the Open method.
	payload []byte

	// The signature of the domain string, type hint, and payload.
	signature []byte
}

var errEmptyDomain = errors.New("envelope domain must not be empty")

// MakeEnvelope constructs a new SignedEnvelope using the given privateKey.
//
// The required 'domain' string contextualizes the envelope for a particular purpose,
// and must be supplied when verifying the signature.
//
// The 'payloadType' field indicates what kind of data is contained and may be empty.
func MakeEnvelope(privateKey PrivKey, domain string, payloadType []byte, contents []byte) (*SignedEnvelope, error) {
	if len(domain) == 0 {
		return nil, errEmptyDomain
	}
	toSign := makeSigBuffer(domain, payloadType, contents)
	sig, err := privateKey.Sign(toSign)
	if err != nil {
		return nil, err
	}

	return &SignedEnvelope{
		PublicKey:   privateKey.GetPublic(),
		PayloadType: payloadType,
		payload:     contents,
		signature:   sig,
	}, nil
}

// UnmarshalEnvelope converts a serialized protobuf representation of an envelope
// into a SignedEnvelope struct.
func UnmarshalEnvelope(serializedEnvelope []byte) (*SignedEnvelope, error) {
	e := pb.SignedEnvelope{}
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
		payload:     e.Payload,
		signature:   e.Signature,
	}, nil
}

// Validate returns true if the envelope signature is valid for the given 'domain',
// or false if it is invalid. May return an error if signature validation fails.
func (e *SignedEnvelope) Validate(domain string) (bool, error) {
	toVerify := makeSigBuffer(domain, e.PayloadType, e.payload)
	return e.PublicKey.Verify(toVerify, e.signature)
}

// Marshal returns a []byte containing a serialized protobuf representation of
// the SignedEnvelope.
func (e *SignedEnvelope) Marshal() ([]byte, error) {
	key, err := PublicKeyToProto(e.PublicKey)
	if err != nil {
		return nil, err
	}
	msg := pb.SignedEnvelope{
		PublicKey:   key,
		PayloadType: e.PayloadType,
		Payload:     e.payload,
		Signature:   e.signature,
	}
	return proto.Marshal(&msg)
}

// Open validates the signature (within the given 'domain') and returns
// the payload of the envelope. Will fail with an error if the signature
// is invalid.
func (e *SignedEnvelope) Open(domain string) ([]byte, error) {
	valid, err := e.Validate(domain)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("invalid signature or incorrect domain")
	}
	return e.payload, nil
}

func (e *SignedEnvelope) Equals(other *SignedEnvelope) bool {
	return e.PublicKey.Equals(other.PublicKey) &&
		bytes.Compare(e.PayloadType, other.PayloadType) == 0 &&
		bytes.Compare(e.payload, other.payload) == 0 &&
		bytes.Compare(e.signature, other.signature) == 0
}

// makeSigBuffer is a helper function that prepares a buffer to sign or verify.
func makeSigBuffer(domain string, typeHint []byte, content []byte) []byte {
	b := bytes.Buffer{}
	domainBytes := []byte(domain)
	b.Write(encodedSize(domainBytes))
	b.Write(domainBytes)
	b.Write(encodedSize(typeHint))
	b.Write(typeHint)
	b.Write(encodedSize(content))
	b.Write(content)
	return b.Bytes()
}

func encodedSize(content []byte) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(len(content)))
	return b
}
