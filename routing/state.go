package routing

import (
	"bytes"
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	pb "github.com/libp2p/go-libp2p-core/routing/pb"
	ma "github.com/multiformats/go-multiaddr"
	"time"
)

// The domain string used for routing state records contained in a SignedEnvelope.
const StateEnvelopeDomain = "libp2p-routing-state"

// The type hint used to identify routing state records in a SignedEnvelope.
// TODO: register multicodec
var StateEnvelopePayloadType = []byte("/libp2p/routing-state-record")

type SignedRoutingState struct {
	// PeerID is the ID of the peer this record pertains to.
	PeerID peer.ID

	// Seq is an increment-only sequence counter used to order RoutingState records in time.
	Seq uint64

	// Addrs contains the public addresses of the peer this record pertains to.
	Addrs []ma.Multiaddr

	// Envelope contains the signature and serialized RoutingStateRecord protobuf.
	// Although it uses a bit
	Envelope *crypto.SignedEnvelope
}

// MakeSignedRoutingState returns a SignedRoutingState record containing the given multiaddrs,
// signed with the given private key.
// It generates a timestamp-based sequence number.
func MakeSignedRoutingState(privKey crypto.PrivKey, addrs []ma.Multiaddr) (*SignedRoutingState, error) {
	p, err := peer.IDFromPrivateKey(privKey)
	if err != nil {
		return nil, err
	}
	idBytes, err := p.MarshalBinary()
	if err != nil {
		return nil, err
	}
	seq := statelessSeqNo()
	msg := pb.RoutingStateRecord{
		PeerId:    idBytes,
		Seq:       seq,
		Addresses: addrsToProtobuf(addrs),
	}
	payload, err := proto.Marshal(&msg)
	if err != nil {
		return nil, err
	}
	envelope, err := crypto.MakeEnvelope(privKey, StateEnvelopeDomain, StateEnvelopePayloadType, payload)
	if err != nil {
		return nil, err
	}
	return &SignedRoutingState{
		PeerID:   p,
		Seq:      seq,
		Addrs:    addrs,
		Envelope: envelope,
	}, nil
}

// UnmarshalSignedRoutingState accepts a serialized SignedEnvelope message containing
// a RoutingStateRecord protobuf and returns a SignedRoutingState record.
// Fails if the signature is invalid, if the envelope has an unexpected payload type,
// if deserialization of the envelope or its inner payload fails.
func UnmarshalSignedRoutingState(envelopeBytes []byte) (*SignedRoutingState, error) {
	envelope, err := crypto.OpenEnvelope(envelopeBytes, StateEnvelopeDomain)
	if err != nil {
		return nil, err
	}
	return SignedRoutingStateFromEnvelope(envelope)
}

// SignedRoutingStateFromEnvelope accepts a SignedEnvelope struct containing
// a RoutingStateRecord protobuf and returns a SignedRoutingState record.
// Fails if the signature is invalid, if the envelope has an unexpected payload type,
// or if deserialization of the envelope payload fails.
func SignedRoutingStateFromEnvelope(envelope *crypto.SignedEnvelope) (*SignedRoutingState, error) {
	if bytes.Compare(envelope.PayloadType, StateEnvelopePayloadType) != 0 {
		return nil, errors.New("unexpected envelope payload type")
	}
	var msg pb.RoutingStateRecord
	err := proto.Unmarshal(envelope.Payload, &msg)
	if err != nil {
		return nil, err
	}
	id, err := peer.IDFromBytes(msg.PeerId)
	if err != nil {
		return nil, err
	}
	if !id.MatchesPublicKey(envelope.PublicKey) {
		return nil, errors.New("peer id in routing state record does not match signing key")
	}
	return &SignedRoutingState{
		PeerID:   id,
		Seq:      msg.Seq,
		Addrs:    addrsFromProtobuf(msg.Addresses),
		Envelope: envelope,
	}, nil
}

// Marshal returns a byte slice containing the SignedRoutingState as a serialized SignedEnvelope
// protobuf message.
func (s *SignedRoutingState) Marshal() ([]byte, error) {
	return s.Envelope.Marshal()
}

// Equal returns true if the other SignedRoutingState is identical to this one.
func (s *SignedRoutingState) Equal(other *SignedRoutingState) bool {
	if other == nil {
		return false
	}
	if s.Seq != other.Seq {
		return false
	}
	if s.PeerID != other.PeerID {
		return false
	}
	if len(s.Addrs) != len(other.Addrs) {
		return false
	}
	for i, _ := range s.Addrs {
		if !s.Addrs[i].Equal(other.Addrs[i]) {
			return false
		}
	}
	return s.Envelope.Equal(other.Envelope)
}

// statelessSeqNo is a helper to generate a timestamp-based sequence number.
func statelessSeqNo() uint64 {
	return uint64(time.Now().UnixNano())
}

func addrsFromProtobuf(addrs []*pb.RoutingStateRecord_AddressInfo) []ma.Multiaddr {
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

func addrsToProtobuf(addrs []ma.Multiaddr) []*pb.RoutingStateRecord_AddressInfo {
	var out []*pb.RoutingStateRecord_AddressInfo
	for _, addr := range addrs {
		out = append(out, &pb.RoutingStateRecord_AddressInfo{Multiaddr: addr.Bytes()})
	}
	return out
}
