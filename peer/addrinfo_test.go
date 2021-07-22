package peer_test

import (
	"testing"

	ma "github.com/multiformats/go-multiaddr"

	. "github.com/libp2p/go-libp2p-core/peer"
)

var (
	testID                         ID
	maddrFull, maddrTpt, maddrPeer ma.Multiaddr
)

func init() {
	var err error
	testID, err = Decode("QmS3zcG7LhYZYSJMhyRZvTddvbNUqtt8BJpaSs6mi1K5Va")
	if err != nil {
		panic(err)
	}
	maddrPeer = ma.StringCast("/p2p/" + Encode(testID))
	maddrTpt = ma.StringCast("/ip4/127.0.0.1/tcp/1234")
	maddrFull = maddrTpt.Encapsulate(maddrPeer)
}

func TestSplitAddr(t *testing.T) {
	tpt, id := SplitAddr(maddrFull)
	if !tpt.Equal(maddrTpt) {
		t.Fatal("expected transport")
	}
	if id != testID {
		t.Fatalf("%s != %s", id, testID)
	}

	tpt, id = SplitAddr(maddrPeer)
	if tpt != nil {
		t.Fatal("expected no transport")
	}
	if id != testID {
		t.Fatalf("%s != %s", id, testID)
	}

	tpt, id = SplitAddr(maddrTpt)
	if !tpt.Equal(maddrTpt) {
		t.Fatal("expected a transport")
	}
	if id != "" {
		t.Fatal("expected no peer ID")
	}
}

func TestAddrInfoFromP2pAddr(t *testing.T) {
	ai, err := AddrInfoFromP2pAddr(maddrFull)
	if err != nil {
		t.Fatal(err)
	}
	if len(ai.Addrs) != 1 || !ai.Addrs[0].Equal(maddrTpt) {
		t.Fatal("expected transport")
	}
	if ai.ID != testID {
		t.Fatalf("%s != %s", ai.ID, testID)
	}

	ai, err = AddrInfoFromP2pAddr(maddrPeer)
	if err != nil {
		t.Fatal(err)
	}
	if len(ai.Addrs) != 0 {
		t.Fatal("expected transport")
	}
	if ai.ID != testID {
		t.Fatalf("%s != %s", ai.ID, testID)
	}

	_, err = AddrInfoFromP2pAddr(maddrTpt)
	if err != ErrInvalidAddr {
		t.Fatalf("wrong error: %s", err)
	}
}

func TestAddrInfosFromP2pAddrs(t *testing.T) {
	infos, err := AddrInfosFromP2pAddrs()
	if err != nil {
		t.Fatal(err)
	}
	if len(infos) != 0 {
		t.Fatal("expected no addrs")
	}
	if _, err = AddrInfosFromP2pAddrs(nil); err == nil {
		t.Fatal("expected nil multiaddr to fail")
	}

	addrs := []ma.Multiaddr{
		ma.StringCast("/ip4/128.199.219.111/tcp/4001/ipfs/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64"),
		ma.StringCast("/ip4/104.236.76.40/tcp/4001/ipfs/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64"),

		ma.StringCast("/ipfs/QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd"),
		ma.StringCast("/ip4/178.62.158.247/tcp/4001/ipfs/QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd"),

		ma.StringCast("/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM"),
	}
	expected := map[string][]ma.Multiaddr{
		"QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64": {
			ma.StringCast("/ip4/128.199.219.111/tcp/4001"),
			ma.StringCast("/ip4/104.236.76.40/tcp/4001"),
		},
		"QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd": {
			ma.StringCast("/ip4/178.62.158.247/tcp/4001"),
		},
		"QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM": nil,
	}
	infos, err = AddrInfosFromP2pAddrs(addrs...)
	if err != nil {
		t.Fatal(err)
	}
	for _, info := range infos {
		exaddrs, ok := expected[info.ID.Pretty()]
		if !ok {
			t.Fatalf("didn't expect peer %s", info.ID)
		}
		if len(info.Addrs) != len(exaddrs) {
			t.Fatalf("got %d addrs, expected %d", len(info.Addrs), len(exaddrs))
		}
		// AddrInfosFromP2pAddrs preserves order. I'd like to keep this
		// guarantee for now.
		for i, addr := range info.Addrs {
			if !exaddrs[i].Equal(addr) {
				t.Fatalf("expected %s, got %s", exaddrs[i], addr)
			}
		}
		delete(expected, info.ID.Pretty())
	}
}

func TestAddrInfoJSON(t *testing.T) {
	ai := AddrInfo{ID: testID, Addrs: []ma.Multiaddr{maddrFull}}
	out, err := ai.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	var addrInfo AddrInfo
	if err := addrInfo.UnmarshalJSON(out); err != nil {
		t.Fatal(err)
	}
	if addrInfo.ID != testID {
		t.Fatalf("expected ID to equal %s, got %s", testID.Pretty(), addrInfo.ID.Pretty())
	}
	if len(addrInfo.Addrs) != 1 || !addrInfo.Addrs[0].Equal(maddrFull) {
		t.Fatalf("expected addrs to match %v, got %v", maddrFull, addrInfo.Addrs)
	}
}
