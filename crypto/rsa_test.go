package crypto

import (
	"math/rand"
	"testing"
)

var (
	testPk PubKey
	testSk PrivKey
)

func init() {
	var err error
	testSk, testPk, err = GenerateRSAKeyPair(2048, rand.New(rand.NewSource(42)))
	if err != nil {
		panic(err)
	}
}

func TestRSABasicSignAndVerify(t *testing.T) {
	data := []byte("hello! and welcome to some awesome crypto primitives")

	sig, err := testSk.Sign(data)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := testPk.Verify(data, sig)
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("signature didnt match")
	}

	// change data
	data[0] = ^data[0]
	ok, err = testPk.Verify(data, sig)
	if err == nil {
		t.Fatal("should have produced a verification error")
	}

	if ok {
		t.Fatal("signature matched and shouldn't")
	}
}

func TestRSASmallKey(t *testing.T) {
	_, _, err := GenerateRSAKeyPair(512, rand.New(rand.NewSource(42)))
	if err != ErrRsaKeyTooSmall {
		t.Fatal("should have refused to create small RSA key")
	}
}

func TestRSASignZero(t *testing.T) {
	data := make([]byte, 0)
	sig, err := testSk.Sign(data)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := testPk.Verify(data, sig)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("signature didn't match")
	}
}

func TestRSAMarshalLoop(t *testing.T) {
	privB, err := testSk.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	privNew, err := UnmarshalPrivateKey(privB)
	if err != nil {
		t.Fatal(err)
	}

	if !testSk.Equals(privNew) || !privNew.Equals(testSk) {
		t.Fatal("keys are not equal")
	}

	pubB, err := testPk.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	pubNew, err := UnmarshalPublicKey(pubB)
	if err != nil {
		t.Fatal(err)
	}

	if !testPk.Equals(pubNew) || !pubNew.Equals(testPk) {
		t.Fatal("keys are not equal")
	}
}
