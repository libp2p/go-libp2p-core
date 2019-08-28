package insecure

import (
	"bytes"
	"context"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/sec"
	"io"
	"net"
	"testing"

	ci "github.com/libp2p/go-libp2p-core/crypto"
)

// Run a set of sessions through the session setup and verification.
func TestConnections(t *testing.T) {
	clientTpt := newTestTransport(t, ci.RSA, 2048)
	serverTpt := newTestTransport(t, ci.Ed25519, 1024)

	testConnection(t, clientTpt, serverTpt)
}

func newTestTransport(t *testing.T, typ, bits int) *Transport {
	priv, pub, err := ci.GenerateKeyPair(typ, bits)
	if err != nil {
		t.Fatal(err)
	}
	id, err := peer.IDFromPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}

	return NewWithIdentity(id, priv)
}

// Create a new pair of connected TCP sockets.
func newConnPair(t *testing.T) (net.Conn, net.Conn) {
	lstnr, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
		return nil, nil
	}

	var clientErr error
	var client net.Conn
	addr := lstnr.Addr()
	done := make(chan struct{})

	go func() {
		defer close(done)
		client, clientErr = net.Dial(addr.Network(), addr.String())
	}()

	server, err := lstnr.Accept()
	<-done

	lstnr.Close()

	if err != nil {
		t.Fatalf("Failed to accept: %v", err)
	}

	if clientErr != nil {
		t.Fatalf("Failed to connect: %v", clientErr)
	}

	return client, server
}

// Create a new pair of connected sessions based off of the provided
// session generators.
func connect(t *testing.T, clientTpt, serverTpt *Transport) (sec.SecureConn, sec.SecureConn) {
	client, server := newConnPair(t)

	// Connect the client and server sessions
	done := make(chan struct{})

	var clientConn sec.SecureConn
	var clientErr error
	go func() {
		defer close(done)
		clientConn, clientErr = clientTpt.SecureOutbound(context.TODO(), client, serverTpt.LocalPeer())
	}()

	serverConn, serverErr := serverTpt.SecureInbound(context.TODO(), server)
	<-done

	if serverErr != nil {
		t.Fatal(serverErr)
	}

	if clientErr != nil {
		t.Fatal(clientErr)
	}

	return clientConn, serverConn
}

// Check the peer IDs
func testIDs(t *testing.T, clientTpt, serverTpt *Transport, clientConn, serverConn sec.SecureConn) {
	if clientConn.LocalPeer() != clientTpt.LocalPeer() {
		t.Fatal("Client Local Peer ID mismatch.")
	}

	if clientConn.RemotePeer() != serverTpt.LocalPeer() {
		t.Fatal("Client Remote Peer ID mismatch.")
	}

	if clientConn.LocalPeer() != serverConn.RemotePeer() {
		t.Fatal("Server Local Peer ID mismatch.")
	}
}

// Check the keys
func testKeys(t *testing.T, clientTpt, serverTpt *Transport, clientConn, serverConn sec.SecureConn) {
	sk := serverConn.LocalPrivateKey()
	pk := sk.GetPublic()

	if !sk.Equals(serverTpt.LocalPrivateKey()) {
		t.Error("Private key Mismatch.")
	}

	if !pk.Equals(clientConn.RemotePublicKey()) {
		t.Error("Public key mismatch.")
	}
}

// Check sending and receiving messages
func testReadWrite(t *testing.T, clientConn, serverConn sec.SecureConn) {
	before := []byte("hello world")
	_, err := clientConn.Write(before)
	if err != nil {
		t.Fatal(err)
	}

	after := make([]byte, len(before))
	_, err = io.ReadFull(serverConn, after)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(before, after) {
		t.Errorf("Message mismatch. %v != %v", before, after)
	}
}

// Setup a new session with a pair of locally connected sockets
func testConnection(t *testing.T, clientTpt, serverTpt *Transport) {
	clientConn, serverConn := connect(t, clientTpt, serverTpt)

	testIDs(t, clientTpt, serverTpt, clientConn, serverConn)
	testKeys(t, clientTpt, serverTpt, clientConn, serverConn)
	testReadWrite(t, clientConn, serverConn)

	clientConn.Close()
	serverConn.Close()
}
