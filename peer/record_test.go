package peer_test

import (
	"github.com/libp2p/go-libp2p-core/crypto"
	. "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/record"
	"github.com/libp2p/go-libp2p-core/test"
	"testing"
)

func TestSignedPeerRecordFromEnvelope(t *testing.T) {
	priv, _, err := test.RandTestKeyPair(crypto.Ed25519, 256)
	test.AssertNilError(t, err)

	addrs := test.GenerateTestAddrs(10)
	id, err := IDFromPrivateKey(priv)
	test.AssertNilError(t, err)

	rec := NewPeerRecord(id, addrs)
	envelope, err := rec.Sign(priv)
	test.AssertNilError(t, err)

	t.Run("is unaltered after round-trip serde", func(t *testing.T) {
		envBytes, err := envelope.Marshal()
		test.AssertNilError(t, err)

		rec2, env2, err := UnmarshalSignedPeerRecord(envBytes)
		test.AssertNilError(t, err)
		if !rec.Equal(rec2) {
			t.Error("expected peer record to be unaltered after round-trip serde")
		}
		if !envelope.Equal(env2) {
			t.Error("expected signed envelope to be unchanged after round-trip serde")
		}
	})

	t.Run("signing fails if signing key does not match peer id in record", func(t *testing.T) {
		id = "some-other-peer-id"
		rec := NewPeerRecord(id, addrs)
		_, err := rec.Sign(priv)
		if err != ErrPeerIdMismatch {
			t.Error("expected signing with mismatched private key to fail")
		}
	})

	t.Run("unwrapping from signed envelope fails if envelope has wrong domain string", func(t *testing.T) {
		payload := []byte("ignored")
		test.AssertNilError(t, err)

		env, err := record.MakeEnvelope(priv, "wrong-domain", PeerRecordEnvelopePayloadType, payload)
		test.AssertNilError(t, err)
		envBytes, err := env.Marshal()
		_, _, err = UnmarshalSignedPeerRecord(envBytes)
		test.ExpectError(t, err, "unwrapping PeerRecord from envelope should fail if envelope was created with wrong domain string")
	})
}
