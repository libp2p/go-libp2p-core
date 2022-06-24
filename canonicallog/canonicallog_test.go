package canonicallog

import (
	"fmt"
	"net"
	"testing"

	"github.com/libp2p/go-libp2p-core/test"
	"github.com/multiformats/go-multiaddr"
)

func TestLogs(t *testing.T) {
	LogMisbehavingPeer(test.RandPeerIDFatal(t), multiaddr.StringCast("/ip4/1.2.3.4"), "somecomponent", fmt.Errorf("something"), "hi")

	netAddr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 80}
	LogMisbehavingPeerNetAddr(test.RandPeerIDFatal(t), netAddr, "somecomponent", fmt.Errorf("something"), "hi")
}
