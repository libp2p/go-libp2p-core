package record

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"time"

	pool "github.com/libp2p/go-buffer-pool"
	"github.com/libp2p/go-libp2p-core/crypto"
	pb "github.com/libp2p/go-libp2p-core/record/pb"

	"github.com/gogo/protobuf/proto"
	"github.com/multiformats/go-varint"
)

// Envelope contains an arbitrary []byte payload, signed by a libp2p peer.
//
// Envelopes are signed in the context of a particular "domain", which is a
// string specified when creating and verifying the envelope. You must know the
// domain string used to produce the envelope in order to verify the signature
// and access the payload.
type Envelope struct {
	// The public key that can be used to verify the signature and derive the peer id of the signer.
	PublicKey crypto.PubKey

	// A binary identifier that indicates what kind of data is contained in the payload.
	// TODO(yusef): enforce multicodec prefix
	PayloadType []byte

	// A monotonically-increasing sequence counter for ordering SignedEnvelopes in time.
	Seq uint64

	// The envelope payload.
	RawPayload []byte

	// The signature of the domain string :: type hint :: payload.
	signature []byte

	// the unmarshalled payload as a Record, cached on first access via the Record accessor method
	cached         Record
	unmarshalError error
	unmarshalOnce  sync.Once
}

var ErrEmptyDomain = errors.New("envelope domain must not be empty")
var ErrEmptyPayloadType = errors.New("payloadType must not be empty")
var ErrInvalidSignature = errors.New("invalid signature or incorrect domain")

// MakeEnvelope constructs a new Envelope using the given privateKey.
//
// The required 'domain' string contextualizes the envelope for a particular purpose,
// and must be supplied when verifying the signature.
//
// The 'PayloadType' field indicates what kind of data is contained and must be non-empty.
func MakeEnvelope(privateKey crypto.PrivKey, domain string, payloadType []byte, payload []byte) (*Envelope, error) {
	if domain == "" {
		return nil, ErrEmptyDomain
	}

	if len(payloadType) == 0 {
		return nil, ErrEmptyPayloadType
	}

	seq := statelessSeqNo()
	unsigned, err := makeUnsigned(domain, payloadType, payload, seq)
	defer pool.Put(unsigned)
	if err != nil {
		return nil, err
	}

	sig, err := privateKey.Sign(unsigned)
	if err != nil {
		return nil, err
	}

	return &Envelope{
		PublicKey:   privateKey.GetPublic(),
		PayloadType: payloadType,
		RawPayload:  payload,
		Seq:         seq,
		signature:   sig,
	}, nil
}

// MakeEnvelopeWithRecord wraps the given Record in an Envelope, and signs it using the given key
// and domain string.
//
// The Record's concrete type must be associated with a payload type identifier
// (see record.RegisterPayloadType).
func MakeEnvelopeWithRecord(privateKey crypto.PrivKey, domain string, rec Record) (*Envelope, error) {
	payloadType, ok := payloadTypeForRecord(rec)
	if !ok {
		return nil, fmt.Errorf("unable to determine value for PayloadType field")
	}
	payloadBytes, err := rec.MarshalRecord()
	if err != nil {
		return nil, fmt.Errorf("error marshaling record: %v", err)
	}
	return MakeEnvelope(privateKey, domain, payloadType, payloadBytes)
}

// ConsumeEnvelope unmarshals a serialized Envelope, and validates its
// signature using the provided 'domain' string. If validation fails, an error
// is returned, along with the unmarshalled envelope so it can be inspected.
// TODO(yusef): improve this doc comment before merge
func ConsumeEnvelope(data []byte, domain string) (envelope *Envelope, rec Record, err error) {
	e, err := UnmarshalEnvelope(data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed when unmarshalling the envelope: %w", err)
	}

	err = e.validate(domain)
	if err != nil {
		return e, nil, fmt.Errorf("failed to validate envelope: %w", err)
	}

	rec, err = e.Record()
	if err != nil {
		return e, nil, fmt.Errorf("failed to unmarshal envelope payload: %w", err)
	}
	return e, rec, nil
}

// TODO(yusef): doc comment before merge
func ConsumeTypedEnvelope(data []byte, domain string, destRecord Record) (envelope *Envelope, err error) {
	e, err := UnmarshalEnvelope(data)
	if err != nil {
		return nil, fmt.Errorf("failed when unmarshalling the envelope: %w", err)
	}

	err = e.validate(domain)
	if err != nil {
		return e, fmt.Errorf("failed to validate envelope: %w", err)
	}

	err = destRecord.UnmarshalRecord(e.RawPayload)
	if err != nil {
		return e, fmt.Errorf("failed to unmarshal envelope payload: %w", err)
	}
	return e, nil
}

// UnmarshalEnvelope unmarshals a serialized Envelope protobuf message,
// without validating its contents. Most users should use ConsumeEnvelope.
func UnmarshalEnvelope(data []byte) (*Envelope, error) {
	var e pb.Envelope
	if err := proto.Unmarshal(data, &e); err != nil {
		return nil, err
	}

	key, err := crypto.PublicKeyFromProto(e.PublicKey)
	if err != nil {
		return nil, err
	}

	return &Envelope{
		PublicKey:   key,
		PayloadType: e.PayloadType,
		RawPayload:  e.Payload,
		Seq:         e.Seq,
		signature:   e.Signature,
	}, nil
}

// Marshal returns a byte slice containing a serialized protobuf representation
// of a Envelope.
func (e *Envelope) Marshal() ([]byte, error) {
	key, err := crypto.PublicKeyToProto(e.PublicKey)
	if err != nil {
		return nil, err
	}

	msg := pb.Envelope{
		PublicKey:   key,
		PayloadType: e.PayloadType,
		Payload:     e.RawPayload,
		Seq:         e.Seq,
		Signature:   e.signature,
	}
	return proto.Marshal(&msg)
}

// Equal returns true if the other Envelope has the same public key,
// payload, payload type, and signature. This implies that they were also
// created with the same domain string.
func (e *Envelope) Equal(other *Envelope) bool {
	if other == nil {
		return e == nil
	}
	return e.Seq == other.Seq &&
		e.PublicKey.Equals(other.PublicKey) &&
		bytes.Compare(e.PayloadType, other.PayloadType) == 0 &&
		bytes.Compare(e.signature, other.signature) == 0 &&
		bytes.Compare(e.RawPayload, other.RawPayload) == 0
}

// Record returns the Envelope's payload unmarshalled as a Record.
// The concrete type of the returned Record depends on which Record
// type was registered for the Envelope's PayloadType - see record.RegisterPayloadType.
//
// Once unmarshalled, the Record is cached for future access.
func (e *Envelope) Record() (Record, error) {
	e.unmarshalOnce.Do(func() {
		if e.cached != nil {
			return
		}
		e.cached, e.unmarshalError = unmarshalRecordPayload(e.PayloadType, e.RawPayload)
	})
	return e.cached, e.unmarshalError
}

// TypedRecord unmarshals the Envelope's payload to the given Record instance.
// It is the caller's responsibility to ensure that the Record type is capable
// of unmarshalling the Envelope payload. Callers can inspect the Envelope's
// PayloadType field to determine the correct type of Record to use.
//
// This method will always unmarshal the Envelope payload even if a cached record
// exists.
func (e *Envelope) TypedRecord(dest Record) error {
	return dest.UnmarshalRecord(e.RawPayload)
}

// validate returns nil if the envelope signature is valid for the given 'domain',
// or an error if signature validation fails.
func (e *Envelope) validate(domain string) error {
	unsigned, err := makeUnsigned(domain, e.PayloadType, e.RawPayload, e.Seq)
	defer pool.Put(unsigned)
	if err != nil {
		return err
	}

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
