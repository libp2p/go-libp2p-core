package crypto

import (
	"crypto/rand"
	mrand "math/rand"
	"testing"
)

func TestRSABasicSignAndVerify(t *testing.T) {
	priv, pub, err := GenerateRSAKeyPair(512, rand.Reader)
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
	_, _, err := GenerateRSAKeyPair(384, rand.Reader)
	if err != ErrRsaKeyTooSmall {
		t.Fatal("should have refused to create small RSA key")
	}
}

func TestRSASignZero(t *testing.T) {
	priv, pub, err := GenerateRSAKeyPair(512, rand.Reader)
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
	priv, pub, err := GenerateRSAKeyPair(512, rand.Reader)
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

func BenchmarkRSAKeyEqualsCheck(b *testing.B) {
	r := mrand.New(mrand.NewSource(42))
	_, k1, _ := GenerateRSAKeyPair(2048, r)
	_, k2, _ := GenerateRSAKeyPair(2048, r)

	for i := 0; i < b.N; i++ {
		if k1.Equals(k2) {
			b.Fatal("this bad")
		}
	}
}
