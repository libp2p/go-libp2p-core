package peer_test

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	ic "github.com/libp2p/go-libp2p-core/crypto"
	. "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/test"

	b58 "github.com/mr-tron/base58/base58"
	mh "github.com/multiformats/go-multihash"
)

var gen1 keyset // generated
var gen2 keyset // generated
var man keyset  // manual

func hash(b []byte) []byte {
	h, _ := mh.Sum(b, mh.SHA2_256, -1)
	return []byte(h)
}

func init() {
	if err := gen1.generate(); err != nil {
		panic(err)
	}
	if err := gen2.generate(); err != nil {
		panic(err)
	}

	skManBytes = strings.Replace(skManBytes, "\n", "", -1)
	if err := man.load(hpkpMan, skManBytes); err != nil {
		panic(err)
	}
}

type keyset struct {
	sk   ic.PrivKey
	pk   ic.PubKey
	hpk  string
	hpkp string
}

func (ks *keyset) generate() error {
	var err error
	ks.sk, ks.pk, err = test.RandTestKeyPair(ic.RSA, 2048)
	if err != nil {
		return err
	}

	bpk, err := ks.pk.Bytes()
	if err != nil {
		return err
	}

	ks.hpk = string(hash(bpk))
	ks.hpkp = b58.Encode([]byte(ks.hpk))
	return nil
}

func (ks *keyset) load(hpkp, skBytesStr string) error {
	skBytes, err := base64.StdEncoding.DecodeString(skBytesStr)
	if err != nil {
		return err
	}

	ks.sk, err = ic.UnmarshalPrivateKey(skBytes)
	if err != nil {
		return err
	}

	ks.pk = ks.sk.GetPublic()
	bpk, err := ks.pk.Bytes()
	if err != nil {
		return err
	}

	ks.hpk = string(hash(bpk))
	ks.hpkp = b58.Encode([]byte(ks.hpk))
	if ks.hpkp != hpkp {
		return fmt.Errorf("hpkp doesn't match key. %s", hpkp)
	}
	return nil
}

func TestIDMatchesPublicKey(t *testing.T) {

	test := func(ks keyset) {
		p1, err := IDB58Decode(ks.hpkp)
		if err != nil {
			t.Fatal(err)
		}

		if ks.hpk != string(p1) {
			t.Error("p1 and hpk differ")
		}

		if !p1.MatchesPublicKey(ks.pk) {
			t.Fatal("p1 does not match pk")
		}

		p2, err := IDFromPublicKey(ks.pk)
		if err != nil {
			t.Fatal(err)
		}

		if p1 != p2 {
			t.Error("p1 and p2 differ", p1.Pretty(), p2.Pretty())
		}

		if p2.Pretty() != ks.hpkp {
			t.Error("hpkp and p2.Pretty differ", ks.hpkp, p2.Pretty())
		}
	}

	test(gen1)
	test(gen2)
	test(man)
}

func TestIDMatchesPrivateKey(t *testing.T) {

	test := func(ks keyset) {
		p1, err := IDB58Decode(ks.hpkp)
		if err != nil {
			t.Fatal(err)
		}

		if ks.hpk != string(p1) {
			t.Error("p1 and hpk differ")
		}

		if !p1.MatchesPrivateKey(ks.sk) {
			t.Fatal("p1 does not match sk")
		}

		p2, err := IDFromPrivateKey(ks.sk)
		if err != nil {
			t.Fatal(err)
		}

		if p1 != p2 {
			t.Error("p1 and p2 differ", p1.Pretty(), p2.Pretty())
		}
	}

	test(gen1)
	test(gen2)
	test(man)
}

func TestPublicKeyExtraction(t *testing.T) {
	t.Skip("disabled until libp2p/go-libp2p-crypto#51 is fixed")
	// Happy path
	_, originalPub, err := ic.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	id, err := IDFromPublicKey(originalPub)
	if err != nil {
		t.Fatal(err)
	}

	extractedPub, err := id.ExtractPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	if extractedPub == nil {
		t.Fatal("failed to extract public key")
	}
	if !originalPub.Equals(extractedPub) {
		t.Fatal("extracted public key doesn't match")
	}

	// Test invalid multihash (invariant of the type of public key)
	pk, err := ID("").ExtractPublicKey()
	if err == nil {
		t.Fatal("expected an error")
	}
	if pk != nil {
		t.Fatal("expected a nil public key")
	}

	// Shouldn't work for, e.g. RSA keys (too large)

	_, rsaPub, err := ic.GenerateKeyPair(ic.RSA, 2048)
	if err != nil {
		t.Fatal(err)
	}
	rsaId, err := IDFromPublicKey(rsaPub)
	if err != nil {
		t.Fatal(err)
	}
	extractedRsaPub, err := rsaId.ExtractPublicKey()
	if err != ErrNoPublicKey {
		t.Fatal(err)
	}
	if extractedRsaPub != nil {
		t.Fatal("expected to fail to extract public key from rsa ID")
	}
}

func TestValidate(t *testing.T) {
	// Empty peer ID invalidates
	err := ID("").Validate()
	if err == nil {
		t.Error("expected error")
	} else if err != ErrEmptyPeerID {
		t.Error("expected error message: " + ErrEmptyPeerID.Error())
	}

	// Non-empty peer ID validates
	p, err := test.RandPeerID()
	if err != nil {
		t.Fatal(err)
	}

	err = p.Validate()
	if err != nil {
		t.Error("expected nil, but found " + err.Error())
	}
}

var hpkpMan = `QmTzw6TAa4QUd4pHhT6rtEbdst2YPkcjPM7uLykAUnPAbM`
var skManBytes = `
CAASpgkwggSiAgEAAoIBAQC7U7Mzuv3KIoOR22rBDc0PSITJEI7Gg6w7Txz6vMutDYl31vY9WPl5Ms/nZgVFQjBUpNppxraCnEDBdkiP++wjWljI8uCam+OlpiMYyqGBDWao0fBL7EAtoSGegognxii+FPCa3jrQOXerw9OdAOiQHizKIkylsEkdEsAkV0Kr50Qasa1e74ofHeaUTcNVFkEDX79zKVKbb9B21b1rs8X58UGaqf+Sd0Hn7r3FJ87sWZUrZ/4zgEtqc7TLcKGzwfgk3IaQT5hrSs73w6qOew+6HjWxAVQJK6iv+E+pV8YKPPaQvxc1ExOS8ggLdrCEfCM/Jz30Ct8VWM8FDyQjw46tAgEDAoIBAHzid3fR/obBrQvnnICz3gowWIYLCdmtHXzfaKcoh8izsPqPTtOQplDMippEA4OBdY3DPEaEeaxoKyukMF/9SBeRkIX3QGcSl8PEF2Xca6teRHCL9Yfy1XPAwRRXBW/ZcH64oGc+0eAmT8fX4mirRbVpczFsMxkgML4MgBg6LHKZCKfmMeerGLSGd8DKv8a97GgRiVC04f77cGEjZfF4JQ/0CrdAmwD0NKUKSVR682xnsU++SzfNz+mkiPZfn46hl7blfeQh1MA2FySZ8QghuCCYhmMXZACmbtnixG1WrlhokJZd0d3IQ4o0Enjbr2RNOqHXhWSPy+bqQCP86TODdxMCgYEA62c7zL7aQtClElNUw3KHLjAaeCawdZ8Yt/oJprMPiIV1YuaYklQgBA9rVnfKKYd1IfT68wrsl1i32j5sgXFJl/Kmau3lw/L8g7B9KSqTnQdMkVnyNFVlNN6T4LIYhKiCo4XMoa9Rth9L7v2KJtrkx9SGqhYKUY5U/6xp7gr3XC8CgYEAy7eclcSUoj+vIJ/I4DhyMDcq+VNpif3d7+sW/c5wBdyNzqEwhLzo7uDz+U6lOEJI6T7PmlTe/DNFDRt0sOyF/J/d3sw3FVf9KJBfl9e03oHs2sJsCbLMzYqQSFM9O90deillurkU+CRX/VijApMjg3fdS8co26ykbyKp2EuG/+MCgYEAnO99Mynm1zXDYYzjLPcEyXVm+sR1o79lz/wGbyIKWwOjl0RltuLAArTyOaUxcQT4wU38ogdIZOXP5tRIVkuGZUxu8fPugqH9rSBTcMcNE1ozC5FMIuOYzem36yFlrcWsbQPdwR+Lzr+H9KkGxJHt2o2vHA6xi7Q4qnLxSVyk6B8CgYEAh8+9uS24bCp0wGqF6tBMICTHUOJGW/6T9UdkqTRKrpMJNGt1rdNF9JX3+4nDetbbRinfvDiUqCIuCLz4dfMD/b/pPzLPY4/+GwrqZTp4lFad5yxIBnczM7G1hYzTfT4TpsZD0dC4pW2P/jsXVwzCV6U+MoTF58htn2xxOt0EqpcCgYAKC3Zz7gK7g9b4Gk971zEWwm69p6bmqKf+5wCLMK3NH7b3ZTM83UiscJbM5kxKVKZ46GN1H7KFAAX2hD5UnKIOHJfFC/70Az2MeyTJDGFmiTry32RU4v97l4YUxhZZtxXEP72+mnJFfc1ezvYZl6fKmImph0aA0aM5rynXHGIjEw==
`
