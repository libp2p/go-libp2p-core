package metrics

import (
	"github.com/libp2p/go-flow-metrics"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

type BandwidthCounter struct {
	totalIn  flow.Meter
	totalOut flow.Meter

	protocolIn  flow.MeterRegistry
	protocolOut flow.MeterRegistry

	peerIn  flow.MeterRegistry
	peerOut flow.MeterRegistry
}

func NewBandwidthCounter() *BandwidthCounter {
	return new(BandwidthCounter)
}

func (bwc *BandwidthCounter) LogSentMessage(size int64) {
	bwc.totalOut.Mark(uint64(size))
}

func (bwc *BandwidthCounter) LogRecvMessage(size int64) {
	bwc.totalIn.Mark(uint64(size))
}

func (bwc *BandwidthCounter) LogSentMessageStream(size int64, proto protocol.ID, p peer.ID) {
	bwc.protocolOut.Get(string(proto)).Mark(uint64(size))
	bwc.peerOut.Get(string(p)).Mark(uint64(size))
}

func (bwc *BandwidthCounter) LogRecvMessageStream(size int64, proto protocol.ID, p peer.ID) {
	bwc.protocolIn.Get(string(proto)).Mark(uint64(size))
	bwc.peerIn.Get(string(p)).Mark(uint64(size))
}

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
