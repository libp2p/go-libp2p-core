package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"io"

	pb "github.com/libp2p/go-libp2p-core/crypto/pb"

	sha256 "github.com/minio/sha256-simd"
)

// ErrRsaKeyTooSmall is returned when trying to generate or parse an RSA key
// that's smaller than 512 bits. Keys need to be larger enough to sign a 256bit
// hash so this is a reasonable absolute minimum.
var ErrRsaKeyTooSmall = errors.New("rsa keys must be >= 512 bits to be useful")

// RsaPrivateKey is an RSA private key.
type RsaPrivateKey struct {
	sk *rsa.PrivateKey
	pk *rsa.PublicKey
}

// RsaPublicKey is an RSA public key.
type RsaPublicKey struct {
	k *rsa.PublicKey
}

// GenerateRSAKeyPair generates a new RSA private and public key.
func GenerateRSAKeyPair(bits int, src io.Reader) (PrivKey, PubKey, error) {
	if bits < 512 {
		return nil, nil, ErrRsaKeyTooSmall
	}
	priv, err := rsa.GenerateKey(src, bits)
	if err != nil {
		return nil, nil, err
	}
	pk := &priv.PublicKey
	return &RsaPrivateKey{sk: priv}, &RsaPublicKey{pk}, nil
}

// Verify compares a signature against input data.
// In the case of RSA keys, verification applies PKCS1v15 RSA signature verification
// upon the SHA-256 hash of the input data.
func (pk *RsaPublicKey) Verify(data, sig []byte) (bool, error) {
	hashed := sha256.Sum256(data)
	err := rsa.VerifyPKCS1v15(pk.k, crypto.SHA256, hashed[:], sig)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (pk *RsaPublicKey) Type() pb.KeyType {
	return pb.KeyType_RSA
}

// Bytes returns protobuf bytes of a public key.
func (pk *RsaPublicKey) Bytes() ([]byte, error) {
	return MarshalPublicKey(pk)
}

// Raw returns a binary representation of the key in a type-specific format.
// In the case of RSA keys, returns an X.509 PKIX encoded public key.
func (pk *RsaPublicKey) Raw() ([]byte, error) {
	return x509.MarshalPKIXPublicKey(pk.k)
}

// Encrypt returns encrypted bytes from the input data.
func (pk *RsaPublicKey) Encrypt(b []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pk.k, b)
}

// Equals checks whether this key is equal to another.
func (pk *RsaPublicKey) Equals(k Key) bool {
	return KeyEqual(pk, k)
}

// Sign returns a signature of the input data.
func (sk *RsaPrivateKey) Sign(message []byte) ([]byte, error) {
	hashed := sha256.Sum256(message)
	return rsa.SignPKCS1v15(rand.Reader, sk.sk, crypto.SHA256, hashed[:])
}

// GetPublic returns a public key.
func (sk *RsaPrivateKey) GetPublic() PubKey {
	if sk.pk == nil {
		sk.pk = &sk.sk.PublicKey
	}
	return &RsaPublicKey{sk.pk}
}

// Decrypt returns decrypted bytes of the input encrypted bytes.
func (sk *RsaPrivateKey) Decrypt(b []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, sk.sk, b)
}

func (sk *RsaPrivateKey) Type() pb.KeyType {
	return pb.KeyType_RSA
}

// Bytes returns protobuf bytes from a private key.
func (sk *RsaPrivateKey) Bytes() ([]byte, error) {
	return MarshalPrivateKey(sk)
}

// Raw returns a binary representation of a private key in a type-specific format.
// In the case of RSA keys, returns the key in ASN.1 DER encoded form.
func (sk *RsaPrivateKey) Raw() ([]byte, error) {
	b := x509.MarshalPKCS1PrivateKey(sk.sk)
	return b, nil
}

// Equals checks whether this key is equal to another.
func (sk *RsaPrivateKey) Equals(k Key) bool {
	return KeyEqual(sk, k)
}

// UnmarshalRsaPrivateKey returns a private key from the input X.509 bytes.
func UnmarshalRsaPrivateKey(b []byte) (PrivKey, error) {
	sk, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, err
	}
	if sk.N.BitLen() < 512 {
		return nil, ErrRsaKeyTooSmall
	}
	return &RsaPrivateKey{sk: sk}, nil
}

// MarshalRsaPrivateKey returns the X.509 bytes of the private key.
func MarshalRsaPrivateKey(k *RsaPrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(k.sk)
}

// UnmarshalRsaPublicKey returns a public key from the input X.509 bytes.
func UnmarshalRsaPublicKey(b []byte) (PubKey, error) {
	pub, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}
	pk, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not actually an rsa public key")
	}
	if pk.N.BitLen() < 512 {
		return nil, ErrRsaKeyTooSmall
	}
	return &RsaPublicKey{pk}, nil
}

// MarshalRsaPublicKey returns the X.509 bytes from the public key.
func MarshalRsaPublicKey(k *RsaPublicKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(k.k)
}
