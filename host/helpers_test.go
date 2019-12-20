package host

import (
	"context"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/routing"
	"github.com/libp2p/go-libp2p-core/test"
	ma "github.com/multiformats/go-multiaddr"
	"testing"
	"time"
)

type mockHost struct {
	fixedPrivKey crypto.PrivKey
	addrs        []ma.Multiaddr
}

func (h *mockHost) Addrs() []ma.Multiaddr {
	return h.addrs
}

func (h *mockHost) Peerstore() peerstore.Peerstore {
	return mockPeerstore{fixedPrivKey: h.fixedPrivKey}
}

func (*mockHost) ID() peer.ID { return "" }

type mockPeerstore struct {
	fixedPrivKey crypto.PrivKey
}

// the one method I care about...
func (m mockPeerstore) PrivKey(peer.ID) crypto.PrivKey {
	return m.fixedPrivKey
}

// so many other things in the Peerstore interface...
func (m mockPeerstore) Close() error                                                { return nil }
func (m mockPeerstore) AddAddr(p peer.ID, addr ma.Multiaddr, ttl time.Duration)     {}
func (m mockPeerstore) AddAddrs(p peer.ID, addrs []ma.Multiaddr, ttl time.Duration) {}
func (m mockPeerstore) AddCertifiedAddrs(s *routing.SignedRoutingState, ttl time.Duration) error {
	return nil
}
func (m mockPeerstore) SetAddr(p peer.ID, addr ma.Multiaddr, ttl time.Duration)           {}
func (m mockPeerstore) SetAddrs(p peer.ID, addrs []ma.Multiaddr, ttl time.Duration)       {}
func (m mockPeerstore) UpdateAddrs(p peer.ID, oldTTL time.Duration, newTTL time.Duration) {}
func (m mockPeerstore) Addrs(p peer.ID) []ma.Multiaddr                                    { return nil }
func (m mockPeerstore) CertifiedAddrs(p peer.ID) []ma.Multiaddr                           { return nil }
func (m mockPeerstore) AddrStream(context.Context, peer.ID) <-chan ma.Multiaddr           { return nil }
func (m mockPeerstore) ClearAddrs(p peer.ID)                                              {}
func (m mockPeerstore) PeersWithAddrs() peer.IDSlice                                      { return nil }
func (m mockPeerstore) SignedRoutingState(p peer.ID) *routing.SignedRoutingState          { return nil }
func (m mockPeerstore) Get(p peer.ID, key string) (interface{}, error)                    { return nil, nil }
func (m mockPeerstore) Put(p peer.ID, key string, val interface{}) error                  { return nil }
func (m mockPeerstore) RecordLatency(peer.ID, time.Duration)                              {}
func (m mockPeerstore) LatencyEWMA(peer.ID) time.Duration                                 { return 0 }
func (m mockPeerstore) GetProtocols(peer.ID) ([]string, error)                            { return nil, nil }
func (m mockPeerstore) AddProtocols(peer.ID, ...string) error                             { return nil }
func (m mockPeerstore) SetProtocols(peer.ID, ...string) error                             { return nil }
func (m mockPeerstore) RemoveProtocols(peer.ID, ...string) error                          { return nil }
func (m mockPeerstore) SupportsProtocols(peer.ID, ...string) ([]string, error)            { return nil, nil }
func (m mockPeerstore) PeerInfo(peer.ID) peer.AddrInfo                                    { return peer.AddrInfo{} }
func (m mockPeerstore) Peers() peer.IDSlice                                               { return nil }
func (m mockPeerstore) PubKey(peer.ID) crypto.PubKey                                      { return nil }
func (m mockPeerstore) AddPubKey(peer.ID, crypto.PubKey) error                            { return nil }
func (m mockPeerstore) AddPrivKey(peer.ID, crypto.PrivKey) error                          { return nil }
func (m mockPeerstore) PeersWithKeys() peer.IDSlice                                       { return nil }

func TestSignedRoutingStateFromHost_FailsIfPrivKeyIsNil(t *testing.T) {
	_, err := SignedRoutingStateFromHost(&mockHost{})
	test.ExpectError(t, err, "expected generating signed routing state to fail when host private key is nil")
}

func TestSignedRoutingStateFromHost_AddrFiltering(t *testing.T) {
	localAddrs := parseAddrs(t,
		// loopback
		"/ip4/127.0.0.1/tcp/42",
		"/ip6/::1/tcp/9999",

		// ip4 LAN reserved
		"/ip4/10.0.0.1/tcp/1234",
		"/ip4/100.64.0.123/udp/10101",
		"/ip4/172.16.0.254/tcp/2345",
		"/ip4/192.168.1.4/udp/1600",

		// link local
		"/ip4/169.254.0.1/udp/1234",
		"/ip6/fe80::c001:37ff:fe6c:0/tcp/42",
	)

	wanAddrs := parseAddrs(t,
		"/ip4/1.2.3.4/tcp/42",
		"/ip4/8.8.8.8/udp/1234",
		"/ip6/2607:f8b0:4002:c02::8a/udp/1234",
		"/ip6/2a03:2880:f111:83:face:b00c:0:25de/udp/2345/quic",
	)

	priv, _, err := test.RandTestKeyPair(crypto.Ed25519, 256)
	if err != nil {
		t.Fatal(err)
	}

	host := &mockHost{
		fixedPrivKey: priv,
		addrs:        append(localAddrs, wanAddrs...),
	}

	// test with local addrs
	state, err := SignedRoutingStateFromHost(host, IncludeLocalAddrs)
	if err != nil {
		t.Fatalf("error generating routing state: %v", err)
	}
	test.AssertAddressesEqual(t, host.addrs, state.Addrs)

	// test filtering out local addrs
	state, err = SignedRoutingStateFromHost(host)
	if err != nil {
		t.Fatalf("error generating routing state: %v", err)
	}
	test.AssertAddressesEqual(t, wanAddrs, state.Addrs)
}

func parseAddrs(t *testing.T, addrStrings ...string) (out []ma.Multiaddr) {
	t.Helper()
	for _, s := range addrStrings {
		out = append(out, ma.StringCast(s))
	}
	return out
}
