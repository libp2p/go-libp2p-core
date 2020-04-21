package test

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc/benchmark/latency"
)

type NetworkTestFunc func(b *testing.B, n1, n2 net.Conn)

func ConnectionForNetwork(n *latency.Network) (n1, n2 net.Conn, err error) {
	var wg sync.WaitGroup
	wg.Add(1)

	var listener net.Listener
	listener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return
	}
	slowListener := n.Listener(listener)
	go func() {
		defer wg.Done()
		n2, _ = slowListener.Accept()
		slowListener.Close()
		return
	}()
	baseDialer := net.Dialer{}
	dialer := n.ContextDialer(baseDialer.DialContext)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	n1, err = dialer(ctx, "tcp4", slowListener.Addr().String())
	if err != nil {
		return
	}

	wg.Wait()
	return
}

func FindNetworkLimit(testFunc NetworkTestFunc) (float64, error) {
	maxProcs := runtime.GOMAXPROCS(0)
	runtime.GOMAXPROCS(maxProcs)
	if maxProcs > 1 {
		runtime.GOMAXPROCS(1)
	}

	network := latency.Network{
		Kbps:    0,
		Latency: 0,
	}

	wrapperFunc := func(sb *testing.B) {
		n1, n2, err := ConnectionForNetwork(&network)
		if err != nil {
			sb.Error(err)
		}
		defer n1.Close()
		defer n2.Close()
		testFunc(sb, n1, n2)
	}

	result := testing.Benchmark(wrapperFunc)
	if result.N < 1 {
		return 0.0, fmt.Errorf("failed to run benchmark")
	}
	max := (float64(result.Bytes) * float64(result.N) / 1e6) / result.T.Seconds()
	fmt.Printf("CPU Bound Limit: %s\n", result)

	current := max
	network.Latency = 500 * time.Microsecond
	for current > max*0.9 {
		network.Latency *= 2
		result = testing.Benchmark(wrapperFunc)
		current = (float64(result.Bytes) * float64(result.N) / 1e6) / result.T.Seconds()
	}
	fmt.Printf("Latency Bound Limit: %s\n", network.Latency)

	network.Kbps = 1024 * 100 // 100Mbps
	network.Latency /= 2
	for current > max*0.9 {
		network.Kbps /= 2
		result = testing.Benchmark(wrapperFunc)
		current = (float64(result.Bytes) * float64(result.N) / 1e6) / result.T.Seconds()
	}
	fmt.Printf("Bandwidth Bound Limit: %dKbps\n", network.Kbps)
	fmt.Printf("Network Bound Limit: %s\n", result)

	// Now, at a network-bounded level (divide bandwidth and latency by 10 as base)
	// look at utilizations as a function of parallelism.
	network.Latency *= 10
	network.Kbps /= 10
	prevBytes := int64(-1)
	ratio := 1.0
	for i := 1; i <= maxProcs; i *= 2 {
		runtime.GOMAXPROCS(i)
		result = testing.Benchmark(wrapperFunc)
		if prevBytes > 0 {
			ratio = float64(result.Bytes) / float64(prevBytes)
		}
		prevBytes = result.Bytes
		fmt.Printf("At MaxProc %d %dKbps / %s latency: %s\n", i, network.Kbps, network.Latency, result)
	}
	fmt.Printf("Slowdown is %f%%\n", 100*(1.0-ratio))
	return (1.0 - ratio), nil
}
