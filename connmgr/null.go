package connmgr

import (
	"context"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

// NullConnMgr is a ConnMgr that provides no functionality. It is used by go-libp2p-blankhost
// to satisfy the ConnMgr interface.
type NullConnMgr struct{}

var _ ConnManager = (*NullConnMgr)(nil)

func (_ NullConnMgr) TagPeer(peer.ID, string, int)      {}
func (_ NullConnMgr) UntagPeer(peer.ID, string)         {}
func (_ NullConnMgr) GetTagInfo(peer.ID) *TagInfo       { return &TagInfo{} }
func (_ NullConnMgr) TrimOpenConns(ctx context.Context) {}
func (_ NullConnMgr) Notifee() network.Notifiee         { return network.GlobalNoopNotifiee }
func (_ NullConnMgr) Protect(peer.ID, string)           {}
func (_ NullConnMgr) Unprotect(peer.ID, string) bool    { return false }
