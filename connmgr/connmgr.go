package connmgr

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"

	ma "github.com/multiformats/go-multiaddr"
)

// ConnManager tracks connections to peers, and allows consumers to associate metadata
// with each peer.
//
// It enables connections to be trimmed based on implementation-defined heuristics.
type ConnManager interface {

	// TagPeer tags a peer with a string, associating a weight with the tag.
	TagPeer(peer.ID, string, int)

	// Untag removes the tagged value from the peer.
	UntagPeer(p peer.ID, tag string)

	// GetTagInfo returns the metadata associated with the peer,
	// or nil if no metadata has been recorded for the peer.
	GetTagInfo(p peer.ID) *TagInfo

	// TrimOpenConns terminates open connections based on an implementation-defined
	// heuristic.
	TrimOpenConns(ctx context.Context)

	// Notifee returns an implementation that can be called back to inform of
	// opened and closed connections.
	Notifee() network.Notifiee

	// Protect protects a peer from having its connection(s) pruned.
	//
	// Tagging allows different parts of the system to manage protections without interfering with one another.
	//
	// Calls to Protect() with the same tag are idempotent. They are not refcounted, so after multiple calls
	// to Protect() with the same tag, a single Unprotect() call bearing the same tag will revoke the protection.
	Protect(id peer.ID, tag string)

	// Unprotect removes a protection that may have been placed on a peer, under the specified tag.
	//
	// The return value indicates whether the peer continues to be protected after this call, by way of a different tag.
	// See notes on Protect() for more info.
	Unprotect(id peer.ID, tag string) (protected bool)
}

// TagInfo stores metadata associated with a peer.
type TagInfo struct {
	FirstSeen time.Time
	Value     int

	// Tags maps tag ids to the numerical values.
	Tags map[string]int

	// Conns maps connection ids (such as remote multiaddr) to their creation time.
	Conns map[string]time.Time
}

type NullConnMgr struct{}

var _ ConnManager = (*NullConnMgr)(nil)

func (_ NullConnMgr) TagPeer(peer.ID, string, int)   {}
func (_ NullConnMgr) UntagPeer(peer.ID, string)      {}
func (_ NullConnMgr) GetTagInfo(peer.ID) *TagInfo    { return &TagInfo{} }
func (_ NullConnMgr) TrimOpenConns(context.Context)  {}
func (_ NullConnMgr) Notifee() network.Notifiee      { return &cmNotifee{} }
func (_ NullConnMgr) Protect(peer.ID, string)        {}
func (_ NullConnMgr) Unprotect(peer.ID, string) bool { return false }

type cmNotifee struct{}

var _ network.Notifiee = (*cmNotifee)(nil)

func (nn *cmNotifee) Connected(n network.Network, c network.Conn)      {}
func (nn *cmNotifee) Disconnected(n network.Network, c network.Conn)   {}
func (nn *cmNotifee) Listen(n network.Network, addr ma.Multiaddr)      {}
func (nn *cmNotifee) ListenClose(n network.Network, addr ma.Multiaddr) {}
func (nn *cmNotifee) OpenedStream(network.Network, network.Stream)     {}
func (nn *cmNotifee) ClosedStream(network.Network, network.Stream)     {}