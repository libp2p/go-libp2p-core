package routing

import (
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/test"
	"testing"
)

func TestSignedRoutingStateFromEnvelope(t *testing.T) {
	priv, _, err := test.RandTestKeyPair(crypto.Ed25519, 256)
	test.AssertNilError(t, err)

	addrs := test.GenerateTestAddrs(10)
	state, err := MakeSignedRoutingState(priv, addrs)
	test.AssertNilError(t, err)

	t.Run("is unaltered after round-trip serde", func(t *testing.T) {
		envBytes, err := state.Marshal()
		test.AssertNilError(t, err)

		state2, err := UnmarshalSignedRoutingState(envBytes)
		test.AssertNilError(t, err)
		if !state.Equal(state2) {
			t.Error("expected routing state to be unaltered after round-trip serde")
		}
	})

	t.Run("unwrapping from signed envelope fails if envelope has wrong domain string", func(t *testing.T) {
		stateBytes, err := state.Marshal()
		test.AssertNilError(t, err)

		env, err := crypto.MakeEnvelope(priv, "wrong-domain", StateEnvelopePayloadType, stateBytes)
		test.AssertNilError(t, err)
		envBytes, err := env.Marshal()
		_, err = UnmarshalSignedRoutingState(envBytes)
		test.ExpectError(t, err, "unwrapping RoutingState from envelope should fail if envelope was created with wrong domain string")
	})
}
