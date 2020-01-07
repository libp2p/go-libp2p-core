package record

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	pool "github.com/libp2p/go-buffer-pool"
	"github.com/libp2p/go-libp2p-core/crypto"
	pb "github.com/libp2p/go-libp2p-core/record/pb"

	"github.com/gogo/protobuf/proto"
	"github.com/multiformats/go-varint"
)

// SignedEnvelope contains an arbitrary []byte payload, signed by a libp2p peer.
//
// Envelopes are signed in the context of a particular "domain", which is a
// string specified when creating and verifying the envelope. You must know the
// domain string used to produce the envelope in order to verify the signature
// and access the payload.
type SignedEnvelope struct {
	// The public key that can be used to verify the signature and derive the peer id of the signer.
	PublicKey crypto.PubKey

	// A binary identifier that indicates what kind of data is contained in the payload.
	// TODO(yusef): enforce multicodec prefix
	PayloadType []byte

	// A monotonically-increasing sequence counter for ordering SignedEnvelopes in time.
	Seq uint64

	// The envelope payload.
	Payload []byte

	// The signature of the domain string :: type hint :: payload.
	signature []byte
}

var ErrEmptyDomain = errors.New("envelope domain must not be empty")
var ErrInvalidSignature = errors.New("invalid signature or incorrect domain")

// MakeEnvelope constructs a new SignedEnvelope using the given privateKey.
//
// The required 'domain' string contextualizes the envelope for a particular purpose,
// and must be supplied when verifying the signature.
//
// The 'PayloadType' field indicates what kind of data is contained and may be empty.
func MakeEnvelope(privateKey crypto.PrivKey, domain string, payloadType []byte, payload []byte) (*SignedEnvelope, error) {
	if domain == "" {
		return nil, ErrEmptyDomain
	}

	seq := statelessSeqNo()
	unsigned, err := makeUnsigned(domain, payloadType, payload, seq)
	if err != nil {
		return nil, err
	}
	defer pool.Put(unsigned)

	sig, err := privateKey.Sign(unsigned)
	if err != nil {
		return nil, err
	}

	return &SignedEnvelope{
		PublicKey:   privateKey.GetPublic(),
		PayloadType: payloadType,
		Payload:     payload,
		Seq:         seq,
		signature:   sig,
	}, nil
}

// ConsumeEnvelope unmarshals a serialized SignedEnvelope, and validates its
// signature using the provided 'domain' string. If validation fails, an error
// is returned, along with the unmarshalled payload so it can be inspected.
func ConsumeEnvelope(data []byte, domain string) (*SignedEnvelope, error) {
	e, err := UnmarshalEnvelope(data)
	if err != nil {
		return nil, fmt.Errorf("failed when unmarshalling the envelope: %w", err)
	}

	return e, e.validate(domain)
}

// UnmarshalEnvelope unmarshals a serialized SignedEnvelope protobuf message,
// without validating its contents. Most users should use ConsumeEnvelope.
func UnmarshalEnvelope(data []byte) (*SignedEnvelope, error) {
	var e pb.SignedEnvelope
	if err := proto.Unmarshal(data, &e); err != nil {
		return nil, err
	}

	key, err := crypto.PublicKeyFromProto(e.PublicKey)
	if err != nil {
		return nil, err
	}

	return &SignedEnvelope{
		PublicKey:   key,
		PayloadType: e.PayloadType,
		Payload:     e.Payload,
		Seq:         e.Seq,
		signature:   e.Signature,
	}, nil
}

// Marshal returns a byte slice containing a serialized protobuf representation
// of a SignedEnvelope.
func (e *SignedEnvelope) Marshal() ([]byte, error) {
	key, err := crypto.PublicKeyToProto(e.PublicKey)
	if err != nil {
		return nil, err
	}

	msg := pb.SignedEnvelope{
		PublicKey:   key,
		PayloadType: e.PayloadType,
		Payload:     e.Payload,
		Seq:         e.Seq,
		Signature:   e.signature,
	}
	return proto.Marshal(&msg)
}

// Equal returns true if the other SignedEnvelope has the same public key,
// payload, payload type, and signature. This implies that they were also
// created with the same domain string.
func (e *SignedEnvelope) Equal(other *SignedEnvelope) bool {
	if other == nil {
		return e == nil
	}
	return e.Seq == other.Seq &&
		e.PublicKey.Equals(other.PublicKey) &&
		bytes.Compare(e.PayloadType, other.PayloadType) == 0 &&
		bytes.Compare(e.signature, other.signature) == 0 &&
		bytes.Compare(e.Payload, other.Payload) == 0
}

// validate returns true if the envelope signature is valid for the given 'domain',
// or false if it is invalid. May return an error if signature validation fails.
func (e *SignedEnvelope) validate(domain string) error {
	unsigned, err := makeUnsigned(domain, e.PayloadType, e.Payload, e.Seq)
	if err != nil {
		return err
	}
	defer pool.Put(unsigned)

	valid, err := e.PublicKey.Verify(unsigned, e.signature)
	if err != nil {
		return fmt.Errorf("failed while verifying signature: %w", err)
	}
	if !valid {
		return ErrInvalidSignature
	}
	return nil
}

// makeUnsigned is a helper function that prepares a buffer to sign or verify.
// It returns a byte slice from a pool. The caller MUST return this slice to the
// pool.
func makeUnsigned(domain string, payloadType []byte, payload []byte, seq uint64) ([]byte, error) {
	var (
		seqBytes = varint.ToUvarint(seq)
		fields   = [][]byte{[]byte(domain), payloadType, seqBytes, payload}

		// fields are prefixed with their length as an unsigned varint. we
		// compute the lengths before allocating the sig buffer so we know how
		// much space to add for the lengths
		flen = make([][]byte, len(fields))
		size = 0
	)

	for i, f := range fields {
		l := len(f)
		flen[i] = varint.ToUvarint(uint64(l))
		size += l + len(flen[i])
	}

	b := pool.Get(size)

	var s int
	for i, f := range fields {
		s += copy(b[s:], flen[i])
		s += copy(b[s:], f)
	}

	return b[:s], nil
}

// statelessSeqNo is a helper to generate a timestamp-based sequence number.
func statelessSeqNo() uint64 {
	return uint64(time.Now().UnixNano())
}
