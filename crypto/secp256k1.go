package crypto

import (
	"fmt"
	"io"

	pb "github.com/libp2p/go-libp2p-core/crypto/pb"
	"github.com/libp2p/go-libp2p-core/internal/catch"

	dcrec "github.com/decred/dcrd/dcrec/secp256k1/v4"
	dcrececdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"github.com/minio/sha256-simd"
)

// Secp256k1PrivateKey is an Secp256k1 private key
type Secp256k1PrivateKey dcrec.PrivateKey

// Secp256k1PublicKey is an Secp256k1 public key
type Secp256k1PublicKey dcrec.PublicKey

// GenerateSecp256k1Key generates a new Secp256k1 private and public key pair
func GenerateSecp256k1Key(src io.Reader) (PrivKey, PubKey, error) {
	privk, err := dcrec.GeneratePrivateKey()
	if err != nil {
		return nil, nil, err
	}

	k := (*Secp256k1PrivateKey)(privk)
	return k, k.GetPublic(), nil
}

// UnmarshalSecp256k1PrivateKey returns a private key from bytes
func UnmarshalSecp256k1PrivateKey(data []byte) (k PrivKey, err error) {
	if len(data) != dcrec.PrivKeyBytesLen {
		return nil, fmt.Errorf("expected secp256k1 data size to be %d", dcrec.PrivKeyBytesLen)
	}
	defer func() { catch.HandlePanic(recover(), &err, "secp256k1 private-key unmarshal") }()

	privk := dcrec.PrivKeyFromBytes(data)
	return (*Secp256k1PrivateKey)(privk), nil
}

// UnmarshalSecp256k1PublicKey returns a public key from bytes
func UnmarshalSecp256k1PublicKey(data []byte) (_k PubKey, err error) {
	defer func() { catch.HandlePanic(recover(), &err, "secp256k1 public-key unmarshal") }()
	k, err := dcrec.ParsePubKey(data)
	if err != nil {
		return nil, err
	}

	return (*Secp256k1PublicKey)(k), nil
}

// Type returns the private key type
func (k *Secp256k1PrivateKey) Type() pb.KeyType {
	return pb.KeyType_Secp256k1
}

// Raw returns the bytes of the key
func (k *Secp256k1PrivateKey) Raw() ([]byte, error) {
	return (*dcrec.PrivateKey)(k).Serialize(), nil
}

// Equals compares two private keys
func (k *Secp256k1PrivateKey) Equals(o Key) bool {
	sk, ok := o.(*Secp256k1PrivateKey)
	if !ok {
		return basicEquals(k, o)
	}

	return k.GetPublic().Equals(sk.GetPublic())
}

// Sign returns a signature from input data
func (k *Secp256k1PrivateKey) Sign(data []byte) (_sig []byte, err error) {
	defer func() { catch.HandlePanic(recover(), &err, "secp256k1 signing") }()
	key := (*dcrec.PrivateKey)(k)
	hash := sha256.Sum256(data)
	sig := dcrececdsa.Sign(key, hash[:])

	return sig.Serialize(), nil
}

// GetPublic returns a public key
func (k *Secp256k1PrivateKey) GetPublic() PubKey {
	return (*Secp256k1PublicKey)((*dcrec.PrivateKey)(k).PubKey())
}

// Type returns the public key type
func (k *Secp256k1PublicKey) Type() pb.KeyType {
	return pb.KeyType_Secp256k1
}

// Raw returns the bytes of the key
func (k *Secp256k1PublicKey) Raw() (res []byte, err error) {
	defer func() { catch.HandlePanic(recover(), &err, "secp256k1 public key marshaling") }()
	return (*dcrec.PublicKey)(k).SerializeCompressed(), nil
}

// Equals compares two public keys
func (k *Secp256k1PublicKey) Equals(o Key) bool {
	sk, ok := o.(*Secp256k1PublicKey)
	if !ok {
		return basicEquals(k, o)
	}

	return (*dcrec.PublicKey)(k).IsEqual((*dcrec.PublicKey)(sk))
}

// Verify compares a signature against the input data
func (k *Secp256k1PublicKey) Verify(data []byte, sigStr []byte) (success bool, err error) {
	defer func() {
		catch.HandlePanic(recover(), &err, "secp256k1 signature verification")

		// To be extra safe.
		if err != nil {
			success = false
		}
	}()
	sig, err := dcrececdsa.ParseDERSignature(sigStr)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(data)
	return sig.Verify(hash[:], (*dcrec.PublicKey)(k)), nil
}
