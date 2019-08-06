package crypto_test

import (
	"bytes"
	"crypto/rand"
	"testing"

	. "github.com/libp2p/go-libp2p-core/crypto"
	pb "github.com/libp2p/go-libp2p-core/crypto/pb"
	"github.com/libp2p/go-libp2p-core/test"
)

func TestKeys(t *testing.T) {
	for _, typ := range KeyTypes {
		testKeyType(typ, t)
	}
}

func testKeyType(typ int, t *testing.T) {
	sk, pk, err := test.RandTestKeyPair(typ, 512)
	if err != nil {
		t.Fatal(err)
	}

	testKeySignature(t, sk)
	testKeyEncoding(t, sk)
	testKeyEquals(t, sk)
	testKeyEquals(t, pk)
}

func testKeySignature(t *testing.T, sk PrivKey) {
	pk := sk.GetPublic()

	text := make([]byte, 16)
	if _, err := rand.Read(text); err != nil {
		t.Fatal(err)
	}

	sig, err := sk.Sign(text)
	if err != nil {
		t.Fatal(err)
	}

	valid, err := pk.Verify(text, sig)
	if err != nil {
		t.Fatal(err)
	}

	if !valid {
		t.Fatal("Invalid signature.")
	}
}

func testKeyEncoding(t *testing.T, sk PrivKey) {
	skbm, err := MarshalPrivateKey(sk)
	if err != nil {
		t.Fatal(err)
	}

	sk2, err := UnmarshalPrivateKey(skbm)
	if err != nil {
		t.Fatal(err)
	}

	if !sk.Equals(sk2) {
		t.Error("Unmarshaled private key didn't match original.\n")
	}

	skbm2, err := MarshalPrivateKey(sk2)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(skbm, skbm2) {
		t.Error("skb -> marshal -> unmarshal -> skb failed.\n", skbm, "\n", skbm2)
	}

	pk := sk.GetPublic()
	pkbm, err := MarshalPublicKey(pk)
	if err != nil {
		t.Fatal(err)
	}

	pk2, err := UnmarshalPublicKey(pkbm)
	if err != nil {
		t.Fatal(err)
	}

	if !pk.Equals(pk2) {
		t.Error("Unmarshaled public key didn't match original.\n")
	}

	pkbm2, err := MarshalPublicKey(pk)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(pkbm, pkbm2) {
		t.Error("skb -> marshal -> unmarshal -> skb failed.\n", pkbm, "\n", pkbm2)
	}
}

func testKeyEquals(t *testing.T, k Key) {
	// kb, err := k.Raw()
	// if err != nil {
	// 	t.Fatal(err)
	// }

	if !KeyEqual(k, k) {
		t.Fatal("Key not equal to itself.")
	}

	// bad test, relies on deep internals..
	// if !KeyEqual(k, testkey(kb)) {
	// 	t.Fatal("Key not equal to key with same bytes.")
	// }

	sk, pk, err := test.RandTestKeyPair(RSA, 512)
	if err != nil {
		t.Fatal(err)
	}

	if KeyEqual(k, sk) {
		t.Fatal("Keys should not equal.")
	}

	if KeyEqual(k, pk) {
		t.Fatal("Keys should not equal.")
	}
}

type testkey []byte

func (pk testkey) Bytes() ([]byte, error) {
	return pk, nil
}

func (pk testkey) Type() pb.KeyType {
	return pb.KeyType_RSA
}

func (pk testkey) Raw() ([]byte, error) {
	return pk, nil
}

func (pk testkey) Equals(k Key) bool {
	if pk.Type() != k.Type() {
		return false
	}
	a, err := pk.Raw()
	if err != nil {
		return false
	}

	b, err := k.Raw()
	if err != nil {
		return false
	}

	return bytes.Equal(a, b)
}

func TestUnknownCurveErrors(t *testing.T) {
	_, _, err := GenerateEKeyPair("P-256")
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = GenerateEKeyPair("error-please")
	if err == nil {
		t.Fatal("expected invalid key type to error")
	}
}

func TestPanicOnUnknownCipherType(t *testing.T) {
	passed := false
	defer func() {
		if !passed {
			t.Fatal("expected known cipher and hash to succeed")
		}
		err := recover()
		errStr, ok := err.(string)
		if !ok {
			t.Fatal("expected string in panic")
		}
		if errStr != "Unrecognized cipher, programmer error?" {
			t.Fatal("expected \"Unrecognized cipher, programmer error?\"")
		}
	}()
	KeyStretcher("AES-256", "SHA1", []byte("foo"))
	passed = true
	KeyStretcher("Fooba", "SHA1", []byte("foo"))
}
