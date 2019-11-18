package routing

import (
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/test"
	"testing"
)

func TestRoutingStateFromAddrInfo(t *testing.T) {
	id, _ := test.RandPeerID()
	addrs := test.GenerateTestAddrs(10)
	info := peer.AddrInfo{
		ID:    id,
		Addrs: addrs,
	}
	state := RoutingStateFromAddrInfo(&info)
	if state.PeerID != info.ID {
		t.Fatalf("expected routing state to have peer id %s, got %s", id.Pretty(), state.PeerID.Pretty())
	}
	test.AssertAddressesEqual(t, addrs, state.Multiaddrs())
}

func TestRoutingStateFromEnvelope(t *testing.T) {
	priv, pub, err := test.RandTestKeyPair(crypto.Ed25519, 256)
	test.AssertNilError(t, err)

	id, err := peer.IDFromPublicKey(pub)
	test.AssertNilError(t, err)

	addrs := test.GenerateTestAddrs(10)
	state := RoutingStateWithMultiaddrs(id, addrs)

	t.Run("can unwrap a RoutingState from a serialized envelope", func(t *testing.T) {
		env, err := state.ToSignedEnvelope(priv)
		test.AssertNilError(t, err)

		envBytes, err := env.Marshal()
		test.AssertNilError(t, err)

		state2, err := RoutingStateFromEnvelope(envBytes)
		if !state.Equal(state2) {
			t.Error("expected routing state to be unaltered after wrapping in signed envelope")
		}
	})

	t.Run("unwrapping from signed envelope fails if peer id does not match signing key", func(t *testing.T) {
		priv2, _, err := test.RandTestKeyPair(crypto.Ed25519, 256)
		test.AssertNilError(t, err)
		env, err := state.ToSignedEnvelope(priv2)
		test.AssertNilError(t, err)
		envBytes, err := env.Marshal()
		test.AssertNilError(t, err)

		_, err = RoutingStateFromEnvelope(envBytes)
		test.ExpectError(t, err, "unwrapping RoutingState from envelope should fail if peer id does not match key used to sign envelope")
	})

	t.Run("unwrapping from signed envelope fails if envelope has wrong domain string", func (t *testing.T) {
		stateBytes, err := state.Marshal()
		test.AssertNilError(t, err)

		env, err := crypto.MakeEnvelope(priv, "wrong-domain", StateEnvelopePayloadType, stateBytes)
		envBytes, err := env.Marshal()
		_, err = RoutingStateFromEnvelope(envBytes)
		test.ExpectError(t, err, "unwrapping RoutingState from envelope should fail if envelope was created with wrong domain string")
	})

	t.Run("unwrapping from signed envelope fails if envelope has wrong payload type", func (t *testing.T) {
		stateBytes, err := state.Marshal()
		test.AssertNilError(t, err)
		payloadType := []byte("wrong-payload-type")
		env, err := crypto.MakeEnvelope(priv, StateEnvelopeDomain, payloadType, stateBytes)
		envBytes, err := env.Marshal()
		_, err = RoutingStateFromEnvelope(envBytes)
		test.ExpectError(t, err, "unwrapping RoutingState from envelope should fail if envelope was created with wrong payload type")
	})
}
