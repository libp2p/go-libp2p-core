package crypto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/golang/protobuf/proto"
	pb "github.com/libp2p/go-libp2p-core/crypto/pb"
)

func MakeEnvelope(privateKey PrivKey, domain string, typeHint []byte, contents []byte) (*pb.SignedEnvelope, error) {
	pubKey, err := PublicKeyToProto(privateKey.GetPublic())
	if err != nil {
		return nil, err
	}

	toSign := makeSigBuffer(domain, typeHint, contents)
	sig, err := privateKey.Sign(toSign)
	if err != nil {
		return nil, err
	}

	return &pb.SignedEnvelope{
		PublicKey: pubKey,
		TypeHint: typeHint,
		Contents: contents,
		Signature: sig,
	}, nil
}

func ValidateEnvelope(domain string, envelope *pb.SignedEnvelope) (bool, error) {
	key, err := PublicKeyFromProto(envelope.PublicKey)
	if err != nil {
		return false, err
	}
	toVerify := makeSigBuffer(domain, envelope.TypeHint, envelope.Contents)
	return key.Verify(toVerify, envelope.Signature)
}

func MarshalEnvelope(envelope *pb.SignedEnvelope) ([]byte, error) {
	return proto.Marshal(envelope)
}

func UnmarshalEnvelope(serializedEnvelope []byte) (*pb.SignedEnvelope, error) {
	e := pb.SignedEnvelope{}
	if err := proto.Unmarshal(serializedEnvelope, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

func OpenEnvelope(domain string, envelope *pb.SignedEnvelope) ([]byte, error) {
	valid, err := ValidateEnvelope(domain, envelope)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("invalid signature")
	}
	return envelope.Contents, nil
}

func makeSigBuffer(domain string, typeHint []byte, content []byte) []byte {
	b := bytes.Buffer{}
	domainBytes := []byte(domain)
	b.Write(encodedSize(domainBytes))
	b.Write(domainBytes)
	b.Write(encodedSize(typeHint))
	b.Write(typeHint)
	b.Write(encodedSize(content))
	b.Write(content)
	return b.Bytes()
}

func encodedSize(content []byte) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(len(content)))
	return b
}