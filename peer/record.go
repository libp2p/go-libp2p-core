package peer

import (
	"errors"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/crypto"
	pb "github.com/libp2p/go-libp2p-core/peer/pb"
	"github.com/libp2p/go-libp2p-core/record"
	ma "github.com/multiformats/go-multiaddr"
)

// The domain string used for peer records contained in a SignedEnvelope.
const PeerRecordEnvelopeDomain = "libp2p-peer-record"

// The type hint used to identify peer records in a SignedEnvelope.
// TODO: register multicodec
var PeerRecordEnvelopePayloadType = []byte("/libp2p/peer-record")

// ErrPeerIdMismatch is returned when attempting to sign a PeerRecord using a key that
// does not match the PeerID contained in the record.
var ErrPeerIdMismatch = errors.New("unable to sign peer record: key does not match record.PeerID")

// PeerRecord contains information that is broadly useful to share with other peers,
// either through a direct exchange (as in the libp2p identify protocol), or through
// a Peer Routing provider, such as a DHT.
//
// Currently, a PeerRecord contains the public listen addresses for a peer, but this
// is expected to expand to include other information in the future.
//
// PeerRecords are intended to be signed by the peer they describe, and there are no
// public methods for marshalling and unmarshalling unsigned PeerRecords.
//
// To share a PeerRecord, first call Sign to wrap the record in a SignedEnvelope
// and sign it with the local peer's private key:
//
//     rec := NewPeerRecord(myPeerId, myAddrs)
//     envelope, err := rec.Sign(myPrivateKey)
//
// The resulting record.SignedEnvelope can be marshalled to a []byte and shared
// publicly. As a convenience, the MarshalSigned method will produce the
// SignedEnvelope and marshal it to a []byte in one go:
//
//     rec := NewPeerRecord(myPeerId, myAddrs)
//     recordBytes, err := rec.MarshalSigned(myPrivateKey)
//
// To validate and unmarshal a signed PeerRecord from a remote peer, use the
// UnmarshalSignedPeerRecord function:
//
//    rec, envelope, err := UnmarshalSignedPeerRecord(recordBytes)
//
// Note that UnmarshalSignedPeerRecord returns the record as well as the
// SignedEnvelope that wraps it, so that you can inspect any metadata
// from the envelope if you need it (for example, the remote peer's public key).
// If you already have an unmarshalled SignedEnvelope, you can call
// PeerRecordFromSignedEnvelope instead:
//
//   rec, err := PeerRecordFromSignedEnvelope(envelope)
type PeerRecord struct {
	// PeerID is the ID of the peer this record pertains to.
	PeerID ID

	// Seq is an increment-only sequence counter used to order peer records in time.
	Seq uint64

	// Addrs contains the public addresses of the peer this record pertains to.
	Addrs []ma.Multiaddr
}

// NewPeerRecord creates a PeerRecord with the given ID and addresses.
// It generates a timestamp-based sequence number.
func NewPeerRecord(id ID, addrs []ma.Multiaddr) *PeerRecord {
	return &PeerRecord{
		PeerID: id,
		Addrs:  addrs,
		Seq:    statelessSeqNo(),
	}
}

// PeerRecordFromAddrInfo creates a PeerRecord from an AddrInfo struct.
// It generates a timestamp-based sequence number.
func PeerRecordFromAddrInfo(info AddrInfo) *PeerRecord {
	return NewPeerRecord(info.ID, info.Addrs)
}

// UnmarshalSignedPeerRecord accepts a []byte containing a SignedEnvelope protobuf message.
// It will try to validate the envelope signature and unmarshal the payload as a PeerRecord.
// Returns the PeerRecord and the SignedEnvelope if successful.
func UnmarshalSignedPeerRecord(envelopeBytes []byte) (*PeerRecord, *record.SignedEnvelope, error) {
	envelope, err := record.ConsumeEnvelope(envelopeBytes, PeerRecordEnvelopeDomain)
	if err != nil {
		return nil, nil, err
	}
	rec, err := PeerRecordFromSignedEnvelope(envelope)
	if err != nil {
		return nil, nil, err
	}
	return rec, envelope, nil
}

// PeerRecordFromSignedEnvelope accepts a SignedEnvelope struct and returns a PeerRecord struct.
// Fails if the signature is invalid, or if the payload cannot be unmarshaled as a PeerRecord.
func PeerRecordFromSignedEnvelope(envelope *record.SignedEnvelope) (*PeerRecord, error) {
	var msg pb.PeerRecord
	err := proto.Unmarshal(envelope.Payload, &msg)
	if err != nil {
		return nil, err
	}
	id, err := IDFromBytes(msg.PeerId)
	if err != nil {
		return nil, err
	}
	if !id.MatchesPublicKey(envelope.PublicKey) {
		return nil, errors.New("peer id in peer record does not match signing key")
	}
	return &PeerRecord{
		PeerID: id,
		Seq:    msg.Seq,
		Addrs:  addrsFromProtobuf(msg.Addresses),
	}, nil
}

// Sign wraps the PeerRecord in a routing.SignedEnvelope, signed with the given
// private key. The private key must match the PeerID field of the PeerRecord.
func (r *PeerRecord) Sign(privKey crypto.PrivKey) (*record.SignedEnvelope, error) {
	p, err := IDFromPrivateKey(privKey)
	if err != nil {
		return nil, err
	}
	if p != r.PeerID {
		return nil, ErrPeerIdMismatch
	}
	idBytes, err := p.MarshalBinary()
	if err != nil {
		return nil, err
	}
	msg := pb.PeerRecord{
		PeerId:    idBytes,
		Seq:       r.Seq,
		Addresses: addrsToProtobuf(r.Addrs),
	}
	payload, err := proto.Marshal(&msg)
	if err != nil {
		return nil, err
	}
	return record.MakeEnvelope(privKey, PeerRecordEnvelopeDomain, PeerRecordEnvelopePayloadType, payload)
}

func (r *PeerRecord) MarshalSigned(privKey crypto.PrivKey) ([]byte, error) {
	env, err := r.Sign(privKey)
	if err != nil {
		return nil, err
	}
	return env.Marshal()
}

// Equal returns true if the other PeerRecord is identical to this one.
func (r *PeerRecord) Equal(other *PeerRecord) bool {
	if other == nil {
		return r == nil
	}
	if r.Seq != other.Seq {
		return false
	}
	if r.PeerID != other.PeerID {
		return false
	}
	if len(r.Addrs) != len(other.Addrs) {
		return false
	}
	for i, _ := range r.Addrs {
		if !r.Addrs[i].Equal(other.Addrs[i]) {
			return false
		}
	}
	return true
}

// statelessSeqNo is a helper to generate a timestamp-based sequence number.
func statelessSeqNo() uint64 {
	return uint64(time.Now().UnixNano())
}

func addrsFromProtobuf(addrs []*pb.PeerRecord_AddressInfo) []ma.Multiaddr {
	var out []ma.Multiaddr
	for _, addr := range addrs {
		a, err := ma.NewMultiaddrBytes(addr.Multiaddr)
		if err != nil {
			continue
		}
		out = append(out, a)
	}
	return out
}

func addrsToProtobuf(addrs []ma.Multiaddr) []*pb.PeerRecord_AddressInfo {
	var out []*pb.PeerRecord_AddressInfo
	for _, addr := range addrs {
		out = append(out, &pb.PeerRecord_AddressInfo{Multiaddr: addr.Bytes()})
	}
	return out
}
