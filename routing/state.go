package routing

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	pb "github.com/libp2p/go-libp2p-core/routing/pb"
)

// The domain string used for routing state records contained in a SignedEnvelope.
const StateEnvelopeDomain = "libp2p-routing-record"

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

// RoutingStateFromEnvelope unwraps a peer RoutingState record from a SignedEnvelope.
// This method will fail if the signature is invalid, or if the record
// belongs to a peer other than the one that signed the envelope.
func RoutingStateFromEnvelope(envelope *crypto.SignedEnvelope) (*RoutingState, error) {
	msgBytes, err := envelope.Open(StateEnvelopeDomain)
	if err != nil {
		return nil, err
	}
	state, err := UnmarshalRoutingState(msgBytes)
	if err != nil {
		return nil, err
	}
	if !state.PeerID.MatchesPublicKey(envelope.PublicKey) {
		return nil, errors.New("peer id in routing state record does not match signing key")
	}
	return state, nil
}

// Marshal serializes a RoutingState record to protobuf and returns its byte representation.
func (s *RoutingState) Marshal() ([]byte, error) {
	id, err := s.PeerID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	msg := pb.RoutingStateRecord{
		PeerId: id,
		Seq: s.Seq,
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