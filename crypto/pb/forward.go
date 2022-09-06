package crypto_pb

import pb "github.com/libp2p/go-libp2p/core/crypto/pb"

var (
	// DEPRECATED use github.com/libp2p/go-libp2p/core/crypto/pb.ErrInvalidLengthCrypto instead
	ErrInvalidLengthCrypto = pb.ErrInvalidLengthCrypto
	// DEPRECATED use github.com/libp2p/go-libp2p/core/crypto/pb.ErrIntOverflowCrypto instead
	ErrIntOverflowCrypto = pb.ErrIntOverflowCrypto
	// DEPRECATED use github.com/libp2p/go-libp2p/core/crypto/pb.ErrUnexpectedEndOfGroupCrypto instead
	ErrUnexpectedEndOfGroupCrypto = pb.ErrUnexpectedEndOfGroupCrypto

	// DEPRECATED use github.com/libp2p/go-libp2p/core/crypto/pb.KeyType_name instead
	KeyType_name = pb.KeyType_name
	// DEPRECATED use github.com/libp2p/go-libp2p/core/crypto/pb.KeyType_value instead
	KeyType_value = pb.KeyType_value
)

// DEPRECATED use github.com/libp2p/go-libp2p/core/crypto/pb.KeyType instead
type KeyType = pb.KeyType

const (
	// DEPRECATED use github.com/libp2p/go-libp2p/core/crypto/pb.KeyType_RSA instead
	KeyType_RSA pb.KeyType = pb.KeyType_RSA
	// DEPRECATED use github.com/libp2p/go-libp2p/core/crypto/pb.KeyType_Ed25519 instead
	KeyType_Ed25519 pb.KeyType = pb.KeyType_Ed25519
	// DEPRECATED use github.com/libp2p/go-libp2p/core/crypto/pb.KeyType_Secp256k1 instead
	KeyType_Secp256k1 pb.KeyType = pb.KeyType_Secp256k1
	// DEPRECATED use github.com/libp2p/go-libp2p/core/crypto/pb.KeyType_ECDSA instead
	KeyType_ECDSA pb.KeyType = pb.KeyType_ECDSA
)

type (
	// DEPRECATED use github.com/libp2p/go-libp2p/core/crypto/pb.PrivateKey instead
	PrivateKey = pb.PrivateKey
	// DEPRECATED use github.com/libp2p/go-libp2p/core/crypto/pb.PublicKey instead
	PublicKey = pb.PublicKey
)
