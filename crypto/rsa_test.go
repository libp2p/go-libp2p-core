package crypto

import (
	"crypto/rand"
	"testing"
)

func TestRSABasicSignAndVerify(t *testing.T) {
	priv, pub, err := GenerateRSAKeyPair(2048, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("hello! and welcome to some awesome crypto primitives")

	sig, err := priv.Sign(data)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := pub.Verify(data, sig)
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("signature didnt match")
	}

	// change data
	data[0] = ^data[0]
	ok, err = pub.Verify(data, sig)
	if err == nil {
		t.Fatal("should have produced a verification error")
	}

	if ok {
		t.Fatal("signature matched and shouldn't")
	}
}

func TestRSASmallKey(t *testing.T) {
	_, _, err := GenerateRSAKeyPair(MinRsaKeyBits/2, rand.Reader)
	if err != ErrRsaKeyTooSmall {
		t.Fatal("should have refused to create small RSA key")
	}
	MinRsaKeyBits /= 2
	badPriv, badPub, err := GenerateRSAKeyPair(MinRsaKeyBits, rand.Reader)
	if err != nil {
		t.Fatalf("should have succeeded, got: %s", err)
	}
	pubBytes, err := MarshalPublicKey(badPub)
	if err != nil {
		t.Fatal(err)
	}
	privBytes, err := MarshalPrivateKey(badPriv)
	if err != nil {
		t.Fatal(err)
	}
	MinRsaKeyBits *= 2
	_, err = UnmarshalPublicKey(pubBytes)
	if err != ErrRsaKeyTooSmall {
		t.Fatal("should have refused to unmarshal a weak key")
	}
	_, err = UnmarshalPrivateKey(privBytes)
	if err != ErrRsaKeyTooSmall {
		t.Fatal("should have refused to unmarshal a weak key")
	}
}

func TestRSASignZero(t *testing.T) {
	priv, pub, err := GenerateRSAKeyPair(2048, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	data := make([]byte, 0)
	sig, err := priv.Sign(data)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := pub.Verify(data, sig)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("signature didn't match")
	}
}

func TestRSAMarshalLoop(t *testing.T) {
	priv, pub, err := GenerateRSAKeyPair(2048, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	privB, err := priv.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	privNew, err := UnmarshalPrivateKey(privB)
	if err != nil {
		t.Fatal(err)
	}

	if !priv.Equals(privNew) || !privNew.Equals(priv) {
		t.Fatal("keys are not equal")
	}

	pubB, err := pub.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	pubNew, err := UnmarshalPublicKey(pubB)
	if err != nil {
		t.Fatal(err)
	}

	if !pub.Equals(pubNew) || !pubNew.Equals(pub) {
		t.Fatal("keys are not equal")
	}
}
