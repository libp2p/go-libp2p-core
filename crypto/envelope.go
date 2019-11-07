package crypto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/golang/protobuf/proto"
	pb "github.com/libp2p/go-libp2p-core/crypto/pb"
	"github.com/libp2p/go-libp2p-core/peer"
)

type SignedEnvelope struct {
	PublicKey PubKey
	TypeHint  []byte
	contents  []byte
	signature []byte
}

func MakeEnvelope(privateKey PrivKey, domain string, typeHint []byte, contents []byte) (*SignedEnvelope, error) {
	toSign := makeSigBuffer(domain, typeHint, contents)
	sig, err := privateKey.Sign(toSign)
	if err != nil {
		return nil, err
	}

	return &SignedEnvelope{
		PublicKey: privateKey.GetPublic(),
		TypeHint:  typeHint,
		contents:  contents,
		signature: sig,
	}, nil
}

func UnmarshalEnvelope(serializedEnvelope []byte) (*SignedEnvelope, error) {
	e := pb.SignedEnvelope{}
	if err := proto.Unmarshal(serializedEnvelope, &e); err != nil {
		return nil, err
	}
	key, err := PublicKeyFromProto(e.PublicKey)
	if err != nil {
		return nil, err
	}
	return &SignedEnvelope{
		PublicKey: key,
		TypeHint:  e.TypeHint,
		contents:  e.Contents,
		signature: e.Signature,
	}, nil
}

func (e *SignedEnvelope) SignerID() (peer.ID, error) {
	return peer.IDFromPublicKey(e.PublicKey)
}

func (e *SignedEnvelope) Validate(domain string) (bool, error) {
	toVerify := makeSigBuffer(domain, e.TypeHint, e.contents)
	return e.PublicKey.Verify(toVerify, e.signature)
}

func (e *SignedEnvelope) Marshal() ([]byte, error) {
	key, err := PublicKeyToProto(e.PublicKey)
	if err != nil {
		return nil, err
	}
	msg := pb.SignedEnvelope{
		PublicKey: key,
		TypeHint: e.TypeHint,
		Contents: e.contents,
		Signature: e.signature,
	}
	return proto.Marshal(&msg)
}

func (e *SignedEnvelope) Open(domain string) ([]byte, error) {
	valid, err := e.Validate(domain)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("invalid signature or incorrect domain")
	}
	return e.contents, nil
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