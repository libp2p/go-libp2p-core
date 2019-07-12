// Package insecure provides an insecure, unencrypted implementation of the the SecureConn and SecureTransport interfaces.
//
// Recommended only for testing and other non-production usage.
package insecure

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/network"
	"net"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/sec"

	ggio "github.com/gogo/protobuf/io"
	ci "github.com/libp2p/go-libp2p-core/crypto"
	pb "github.com/libp2p/go-libp2p-core/sec/insecure/pb"
)

// ID is the multistream-select protocol ID that should be used when identifying
// this security transport.
const ID = "/plaintext/2.0.0"

// Transport is a no-op stream security transport. It provides no
// security and simply mocks the security methods. Identity methods
// return the local peer's ID and private key, and whatever the remote
// peer presents as their ID and public key.
// No authentication of the remote identity is performed.
type Transport struct {
	id  peer.ID
	key ci.PrivKey
}

// New constructs a new insecure transport.
func New(id peer.ID, key ci.PrivKey) *Transport {
	return &Transport{
		id:  id,
		key: key,
	}
}

// LocalPeer returns the transport's local peer ID.
func (t *Transport) LocalPeer() peer.ID {
	return t.id
}

// LocalPrivateKey returns the local private key.
// This key is used only for identity generation and provides no security.
func (t *Transport) LocalPrivateKey() ci.PrivKey {
	return t.key
}

// SecureInbound *pretends to secure* an outbound connection to the given peer.
// It sends the local peer's ID and public key, and receives the same from the remote peer.
// No validation is performed as to the authenticity or ownership of the provided public key,
// and the key exchange provides no security.
//
// SecureInbound may fail if the remote peer sends an ID and public key that are inconsistent
// with each other, or if a network error occurs during the ID exchange.
func (t *Transport) SecureInbound(ctx context.Context, insecure net.Conn) (sec.SecureConn, error) {
	conn := &Conn{
		Conn:         insecure,
		local:        t.id,
		localPrivKey: t.key,
	}

	err := conn.runHandshakeSync(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// SecureOutbound *pretends to secure* an outbound connection to the given peer.
// It sends the local peer's ID and public key, and receives the same from the remote peer.
// No validation is performed as to the authenticity or ownership of the provided public key,
// and the key exchange provides no security.
//
// SecureOutbound may fail if the remote peer sends an ID and public key that are inconsistent
// with each other, or if the ID sent by the remote peer does not match the one dialed. It may
// also fail if a network error occurs during the ID exchange.
func (t *Transport) SecureOutbound(ctx context.Context, insecure net.Conn, p peer.ID) (sec.SecureConn, error) {
	conn := &Conn{
		Conn:         insecure,
		local:        t.id,
		localPrivKey: t.key,
	}

	err := conn.runHandshakeSync(ctx)
	if err != nil {
		return nil, err
	}

	if p != conn.remote {
		return nil, fmt.Errorf("remote peer sent unexpected peer ID. expected=%s received=%s",
			p, conn.remote)
	}

	return conn, nil
}

// Conn is the connection type returned by the insecure transport.
type Conn struct {
	net.Conn

	local  peer.ID
	remote peer.ID

	localPrivKey ci.PrivKey
	remotePubKey ci.PubKey
}

func makeExchangeMessage(privkey ci.PrivKey) (*pb.Exchange, error) {
	pubkey, err := ci.PublicKeyToProto(privkey.GetPublic())
	if err != nil {
		return nil, err
	}
	id, err := peer.IDFromPrivateKey(privkey)
	if err != nil {
		return nil, err
	}

	return &pb.Exchange{
		Id:     []byte(id),
		Pubkey: pubkey,
	}, nil
}

func (ic *Conn) runHandshakeSync(ctx context.Context) error {
	reader := ggio.NewDelimitedReader(ic.Conn, network.MessageSizeMax)
	writer := ggio.NewDelimitedWriter(ic.Conn)

	// Generate an Exchange message
	msg, err := makeExchangeMessage(ic.localPrivKey)
	if err != nil {
		return err
	}

	// Send our Exchange and read theirs
	err = writer.WriteMsg(msg)
	if err != nil {
		return err
	}

	remoteMsg := new(pb.Exchange)
	err = reader.ReadMsg(remoteMsg)
	if err != nil {
		return err
	}

	// Pull remote ID and public key from message
	remotePubkey, err := ci.PublicKeyFromProto(*remoteMsg.Pubkey)
	if err != nil {
		return err
	}

	remoteID, err := peer.IDFromPublicKey(remotePubkey)
	if err != nil {
		return err
	}

	// Validate that ID matches public key
	if !remoteID.MatchesPublicKey(remotePubkey) {
		calculatedID, _ := peer.IDFromPublicKey(remotePubkey)
		return fmt.Errorf("remote peer id does not match public key. id=%s calculated_id=%s",
			remoteID, calculatedID)
	}

	// Add remote ID and key to conn state
	ic.remotePubKey = remotePubkey
	ic.remote = remoteID
	return nil
}

// LocalPeer returns the local peer ID.
func (ic *Conn) LocalPeer() peer.ID {
	return ic.local
}

// RemotePeer returns the remote peer ID if we initiated the dial. Otherwise, it
// returns "" (because this connection isn't actually secure).
func (ic *Conn) RemotePeer() peer.ID {
	return ic.remote
}

// RemotePublicKey returns whatever public key was given by the remote peer.
// Note that no verification of ownership is done, as this connection is not secure.
func (ic *Conn) RemotePublicKey() ci.PubKey {
	return ic.remotePubKey
}

// LocalPrivateKey returns the private key for the local peer.
func (ic *Conn) LocalPrivateKey() ci.PrivKey {
	return ic.localPrivKey
}

var _ sec.SecureTransport = (*Transport)(nil)
var _ sec.SecureConn = (*Conn)(nil)
