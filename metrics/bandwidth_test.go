package metrics

import (
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

func BenchmarkBandwidthCounter(b *testing.B) {
	b.StopTimer()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bwc := NewBandwidthCounter()
		round(bwc, b)
	}
}

func round(bwc *BandwidthCounter, b *testing.B) {
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(10000)
	for i := 0; i < 1000; i++ {
		p := peer.ID(fmt.Sprintf("peer-%d", i))
		for j := 0; j < 10; j++ {
			proto := protocol.ID(fmt.Sprintf("bitswap-%d", j))
			go func() {
				defer wg.Done()
				<-start

				for i := 0; i < 1000; i++ {
					bwc.LogSentMessage(100)
					bwc.LogSentMessageStream(100, proto, p)
					time.Sleep(1 * time.Millisecond)
				}
			}()
		}
	}

	b.StartTimer()
	close(start)
	wg.Wait()
	b.StopTimer()
}

// Allow 1% errors for bw calculations.
const acceptableError = 0.01

func TestBandwidthCounter(t *testing.T) {
	bwc := NewBandwidthCounter()
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(200)
	for i := 0; i < 100; i++ {
		p := peer.ID(fmt.Sprintf("peer-%d", i))
		for j := 0; j < 2; j++ {
			proto := protocol.ID(fmt.Sprintf("proto-%d", j))
			go func() {
				defer wg.Done()

				// make sure the bandwidth counters are active
				bwc.LogSentMessage(100)
				bwc.LogRecvMessage(50)
				bwc.LogSentMessageStream(100, proto, p)
				bwc.LogRecvMessageStream(50, proto, p)

				<-start

				t := time.NewTicker(100 * time.Millisecond)
				defer t.Stop()

				for i := 0; i < 39; i++ {
					bwc.LogSentMessage(100)
					bwc.LogRecvMessage(50)
					bwc.LogSentMessageStream(100, proto, p)
					bwc.LogRecvMessageStream(50, proto, p)
					<-t.C
				}
			}()
		}
	}

	assertProtocols := func(check func(Stats)) {
		byProtocol := bwc.GetBandwidthByProtocol()
		if len(byProtocol) != 2 {
			t.Errorf("expected 2 protocols, got %d", len(byProtocol))
		}
		for i := 0; i < 2; i++ {
			p := protocol.ID(fmt.Sprintf("proto-%d", i))
			for _, stats := range [...]Stats{bwc.GetBandwidthForProtocol(p), byProtocol[p]} {
				check(stats)
			}
		}
	}

	assertPeers := func(check func(Stats)) {
		byPeer := bwc.GetBandwidthByPeer()
		if len(byPeer) != 100 {
			t.Errorf("expected 100 peers, got %d", len(byPeer))
		}
		for i := 0; i < 100; i++ {
			p := peer.ID(fmt.Sprintf("peer-%d", i))
			for _, stats := range [...]Stats{bwc.GetBandwidthForPeer(p), byPeer[p]} {
				check(stats)
			}
		}
	}

	time.Sleep(time.Second)
	close(start)
	time.Sleep(2*time.Second + 100*time.Millisecond)

	assertPeers(func(stats Stats) {
		assertApproxEq(t, 2000, stats.RateOut)
		assertApproxEq(t, 1000, stats.RateIn)
	})

	assertProtocols(func(stats Stats) {
		assertApproxEq(t, 100000, stats.RateOut)
		assertApproxEq(t, 50000, stats.RateIn)
	})

	{
		stats := bwc.GetBandwidthTotals()
		assertApproxEq(t, 200000, stats.RateOut)
		assertApproxEq(t, 100000, stats.RateIn)
	}

	wg.Wait()
	time.Sleep(1 * time.Second)

	assertPeers(func(stats Stats) {
		assertEq(t, 8000, stats.TotalOut)
		assertEq(t, 4000, stats.TotalIn)
	})

	assertProtocols(func(stats Stats) {
		assertEq(t, 400000, stats.TotalOut)
		assertEq(t, 200000, stats.TotalIn)
	})

	{
		stats := bwc.GetBandwidthTotals()
		assertEq(t, 800000, stats.TotalOut)
		assertEq(t, 400000, stats.TotalIn)
	}
}

func TestResetBandwidthCounter(t *testing.T) {
	bwc := NewBandwidthCounter()

	p := peer.ID("peer-0")
	proto := protocol.ID("proto-0")

	// We don't calculate bandwidth till we've been active for a second.
	bwc.LogSentMessage(42)
	bwc.LogRecvMessage(24)
	bwc.LogSentMessageStream(100, proto, p)
	bwc.LogRecvMessageStream(50, proto, p)

	time.Sleep(1*time.Second + time.Millisecond)

	bwc.LogSentMessage(42)
	bwc.LogRecvMessage(24)
	bwc.LogSentMessageStream(100, proto, p)
	bwc.LogRecvMessageStream(50, proto, p)

	time.Sleep(1*time.Second + time.Millisecond)

	{
		stats := bwc.GetBandwidthTotals()
		assertEq(t, 84, stats.TotalOut)
		assertEq(t, 48, stats.TotalIn)
	}

	{
		stats := bwc.GetBandwidthByProtocol()
		assertApproxEq(t, 1, float64(len(stats)))
		stat := stats[proto]
		assertApproxEq(t, 100, stat.RateOut)
		assertApproxEq(t, 50, stat.RateIn)
	}

	{
		stats := bwc.GetBandwidthByPeer()
		assertApproxEq(t, 1, float64(len(stats)))
		stat := stats[p]
		assertApproxEq(t, 100, stat.RateOut)
		assertApproxEq(t, 50, stat.RateIn)
	}

	bwc.Reset()
	{
		stats := bwc.GetBandwidthTotals()
		assertEq(t, 0, stats.TotalOut)
		assertEq(t, 0, stats.TotalIn)
	}

	{
		byProtocol := bwc.GetBandwidthByProtocol()
		if len(byProtocol) != 0 {
			t.Errorf("expected 0 protocols, got %d", len(byProtocol))
		}
	}

	{
		byPeer := bwc.GetBandwidthByPeer()
		if len(byPeer) != 0 {
			t.Errorf("expected 0 peers, got %d", len(byPeer))
		}
	}
}

func assertEq(t *testing.T, expected, actual int64) {
	if expected != actual {
		t.Errorf("expected  %d, got %d", expected, actual)
	}
}

func assertApproxEq(t *testing.T, expected, actual float64) {
	t.Helper()
	margin := expected * acceptableError
	if !(math.Abs(expected-actual) <= margin) {
		t.Errorf("expected %f (±%f), got %f", expected, margin, actual)
	}
}
