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

	rec := &PeerRecord{PeerID: id, Addrs: addrs}
	envelope, err := rec.Sign(priv)
	test.AssertNilError(t, err)

	//t.Run("sanity check, don't push to remote", func(t *testing.T) {
	//	id.UnmarshalBinary()
	//})

	t.Run("is unaltered after round-trip serde", func(t *testing.T) {
		envBytes, err := envelope.Marshal()
		test.AssertNilError(t, err)

		env2, untypedRecord, err := record.ConsumeEnvelope(envBytes, PeerRecordEnvelopeDomain)
		test.AssertNilError(t, err)
		rec2, ok := untypedRecord.(*PeerRecord)
		if !ok {
			t.Error("unmarshaled record is not a *PeerRecord")
		}
		if !rec.Equal(rec2) {
			t.Error("expected peer record to be unaltered after round-trip serde")
		}
		if !envelope.Equal(env2) {
			t.Error("expected signed envelope to be unchanged after round-trip serde")
		}
	})

	t.Run("signing fails if signing key does not match peer id in record", func(t *testing.T) {
		id = "some-other-peer-id"
		rec := &PeerRecord{PeerID: id, Addrs: addrs}
		_, err := rec.Sign(priv)
		if err != ErrPeerIdMismatch {
			t.Error("expected signing with mismatched private key to fail")
		}
	})
}
