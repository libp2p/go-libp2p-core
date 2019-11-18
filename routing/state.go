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

// AnnotatedAddr will extend the Multiaddr type with additional metadata, as
// extensions are added to the routing state record spec. It's defined now to
// make refactoring simpler in the future.
type AnnotatedAddr struct {
	ma.Multiaddr
}

// RoutingState contains a snapshot of public, transient state (e.g. addresses, supported protocols)
// for a peer at a given point in time, where "time" is defined by the sequence counter
// field Seq. Greater Seq values are later in time than lesser values, but there are no
// guarantees about the wall-clock time between any two Seq values.
//
// Note that Seq values are peer-specific and can only be compared for records with equal PeerIDs.
type RoutingState struct {
	// PeerID is the ID of the peer this record pertains to.
	PeerID peer.ID

	// Seq is an increment-only sequence counter used to order RoutingState records in time.
	Seq uint64

	// Addresses contains the public addresses of the peer this record pertains to.
	Addresses []*AnnotatedAddr
}

// RoutingStateWithMultiaddrs returns a RoutingState record for the given peer id
// that contains the given multiaddrs. It generates a timestamp-based sequence number.
func RoutingStateWithMultiaddrs(p peer.ID, addrs []ma.Multiaddr) *RoutingState {
	annotated := make([]*AnnotatedAddr, len(addrs))
	for i, a := range addrs {
		annotated[i] = &AnnotatedAddr{Multiaddr: a}
	}
	return &RoutingState{
		PeerID:    p,
		Seq:       statelessSeqNo(),
		Addresses: annotated,
	}
}

// RoutingStateFromAddrInfo converts a peer.AddrInfo into a RoutingState record.
// It generates a timestamp-based sequence number.
func RoutingStateFromAddrInfo(info *peer.AddrInfo) *RoutingState {
	return RoutingStateWithMultiaddrs(info.ID, info.Addrs)
}

// UnmarshalRoutingState unpacks a peer RoutingState record from a serialized protobuf representation.
func UnmarshalRoutingState(serialized []byte) (*RoutingState, error) {
	msg := pb.RoutingStateRecord{}
	err := proto.Unmarshal(serialized, &msg)
	if err != nil {
		return nil, err
	}
	id, err := peer.IDFromBytes(msg.PeerId)
	if err != nil {
		return nil, err
	}
	return &RoutingState{
		PeerID:    id,
		Seq:       msg.Seq,
		Addresses: addrsFromProtobuf(msg.Addresses),
	}, nil
}

// RoutingStateFromEnvelope unwraps a peer RoutingState record from a serialized SignedEnvelope.
// This method will fail if the signature is invalid, or if the record
// belongs to a peer other than the one that signed the envelope.
func RoutingStateFromEnvelope(envelopeBytes []byte) (*RoutingState, error) {
	envelope, err := crypto.OpenEnvelope(envelopeBytes, StateEnvelopeDomain)
	if err != nil {
		return nil, err
	}
	if bytes.Compare(envelope.PayloadType(), StateEnvelopePayloadType) != 0 {
		return nil, errors.New("unexpected envelope payload type")
	}
	state, err := UnmarshalRoutingState(envelope.Payload())
	if err != nil {
		return nil, err
	}
	if !state.PeerID.MatchesPublicKey(envelope.PublicKey()) {
		return nil, errors.New("peer id in routing state record does not match signing key")
	}
	return state, nil
}

// ToSignedEnvelope wraps a Marshal'd RoutingState record in a SignedEnvelope using the
// given private signing key.
func (s *RoutingState) ToSignedEnvelope(key crypto.PrivKey) (*crypto.SignedEnvelope, error) {
	payload, err := s.Marshal()
	if err != nil {
		return nil, err
	}
	return crypto.MakeEnvelope(key, StateEnvelopeDomain, StateEnvelopePayloadType, payload)
}

// Marshal serializes a RoutingState record to protobuf and returns its byte representation.
func (s *RoutingState) Marshal() ([]byte, error) {
	id, err := s.PeerID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	msg := pb.RoutingStateRecord{
		PeerId:    id,
		Seq:       s.Seq,
		Addresses: addrsToProtobuf(s.Addresses),
	}
	return proto.Marshal(&msg)
}

// Multiaddrs returns the addresses from a RoutingState record without any metadata annotations.
func (s *RoutingState) Multiaddrs() []ma.Multiaddr {
	out := make([]ma.Multiaddr, len(s.Addresses))
	for i, addr := range s.Addresses {
		out[i] = addr.Multiaddr
	}
	return out
}

func (s *RoutingState) Equal(other *RoutingState) bool {
	if s.Seq != other.Seq {
		return false
	}
	if s.PeerID != other.PeerID {
		return false
	}
	if len(s.Addresses) != len(other.Addresses) {
		return false
	}
	for i, _ := range s.Addresses {
		if !s.Addresses[i].Equal(other.Addresses[i]) {
			return false
		}
	}
	return true
}

func (a *AnnotatedAddr) Equal(other *AnnotatedAddr) bool {
	return a.Multiaddr.Equal(other.Multiaddr)
}

func statelessSeqNo() uint64 {
	return uint64(time.Now().UnixNano())
}

func addrsFromProtobuf(addrs []*pb.RoutingStateRecord_AddressInfo) []*AnnotatedAddr {
	out := make([]*AnnotatedAddr, 0)
	for _, addr := range addrs {
		a, err := ma.NewMultiaddrBytes(addr.Multiaddr)
		if err != nil {
			continue
		}
		out = append(out, &AnnotatedAddr{Multiaddr: a})
	}
	return out
}

func addrsToProtobuf(addrs []*AnnotatedAddr) []*pb.RoutingStateRecord_AddressInfo {
	out := make([]*pb.RoutingStateRecord_AddressInfo, 0)
	for _, addr := range addrs {
		out = append(out, &pb.RoutingStateRecord_AddressInfo{Multiaddr: addr.Bytes()})
	}
	return out
}
