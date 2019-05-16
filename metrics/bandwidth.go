// Package metrics provides metrics collection and reporting interfaces for libp2p.
package metrics

import (
	"github.com/libp2p/go-flow-metrics"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

// BandwidthCounter tracks incoming and outgoing data transferred by the local peer.
// Metrics are available for total bandwidth across all peers / protocols, as well
// as segmented by remote peer ID and protocol ID.
type BandwidthCounter struct {
	totalIn  flow.Meter
	totalOut flow.Meter

	protocolIn  flow.MeterRegistry
	protocolOut flow.MeterRegistry

	peerIn  flow.MeterRegistry
	peerOut flow.MeterRegistry
}

// NewBandwidthCounter creates a new BandwidthCounter.
func NewBandwidthCounter() *BandwidthCounter {
	return new(BandwidthCounter)
}

// LogSentMessage records the size of an outgoing message
// without associating the bandwidth to a specific peer or protocol.
func (bwc *BandwidthCounter) LogSentMessage(size int64) {
	bwc.totalOut.Mark(uint64(size))
}

// LogRecvMessage records the size of an incoming message
// without associating the bandwith to a specific peer or protocol.
func (bwc *BandwidthCounter) LogRecvMessage(size int64) {
	bwc.totalIn.Mark(uint64(size))
}

// LogSentMessageStream records the size of an outgoing message over a single logical stream.
// Bandwidth is associated with the given protocol.ID and peer.ID.
func (bwc *BandwidthCounter) LogSentMessageStream(size int64, proto protocol.ID, p peer.ID) {
	bwc.protocolOut.Get(string(proto)).Mark(uint64(size))
	bwc.peerOut.Get(string(p)).Mark(uint64(size))
}

// LogRecvMessageStream records the size of an incoming message over a single logical stream.
// Bandwidth is associated with the given protocol.ID and peer.ID.
func (bwc *BandwidthCounter) LogRecvMessageStream(size int64, proto protocol.ID, p peer.ID) {
	bwc.protocolIn.Get(string(proto)).Mark(uint64(size))
	bwc.peerIn.Get(string(p)).Mark(uint64(size))
}

// GetBandwidthForPeer returns a Stats struct with bandwidth metrics associated with the given peer.ID.
// The metrics returned include all traffic sent / received for the peer, regardless of protocol.
func (bwc *BandwidthCounter) GetBandwidthForPeer(p peer.ID) (out Stats) {
	inSnap := bwc.peerIn.Get(string(p)).Snapshot()
	outSnap := bwc.peerOut.Get(string(p)).Snapshot()

	return Stats{
		TotalIn:  int64(inSnap.Total),
		TotalOut: int64(outSnap.Total),
		RateIn:   inSnap.Rate,
		RateOut:  outSnap.Rate,
	}
}

// GetBandwidthForProtocol returns a Stats struct with bandwidth metrics associated with the given protocol.ID.
// The metrics returned include all traffic sent / recieved for the protocol, regardless of which peers were
// involved.
func (bwc *BandwidthCounter) GetBandwidthForProtocol(proto protocol.ID) (out Stats) {
	inSnap := bwc.protocolIn.Get(string(proto)).Snapshot()
	outSnap := bwc.protocolOut.Get(string(proto)).Snapshot()

	return Stats{
		TotalIn:  int64(inSnap.Total),
		TotalOut: int64(outSnap.Total),
		RateIn:   inSnap.Rate,
		RateOut:  outSnap.Rate,
	}
}

// GetBandwidthTotals returns a Stats struct with bandwidth metrics for all data sent / recieved by the
// local peer, regardless of protocol or remote peer IDs.
func (bwc *BandwidthCounter) GetBandwidthTotals() (out Stats) {
	inSnap := bwc.totalIn.Snapshot()
	outSnap := bwc.totalOut.Snapshot()

	return Stats{
		TotalIn:  int64(inSnap.Total),
		TotalOut: int64(outSnap.Total),
		RateIn:   inSnap.Rate,
		RateOut:  outSnap.Rate,
	}
}
