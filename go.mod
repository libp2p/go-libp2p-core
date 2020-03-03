module github.com/libp2p/go-libp2p-core

require (
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/coreos/go-semver v0.3.0
	github.com/gogo/protobuf v1.3.1
	github.com/ipfs/go-cid v0.0.5
	github.com/jbenet/goprocess v0.1.3
	github.com/libp2p/go-buffer-pool v0.0.2
	github.com/libp2p/go-flow-metrics v0.0.3
	github.com/libp2p/go-openssl v0.0.4
	github.com/minio/sha256-simd v0.1.1
	github.com/mr-tron/base58 v1.1.3
	github.com/multiformats/go-multiaddr v0.2.0
	github.com/multiformats/go-multihash v0.0.13
	github.com/multiformats/go-varint v0.0.5
	github.com/smola/gocompat v0.2.0
	go.opencensus.io v0.22.3
)

go 1.13

replace github.com/libp2p/go-addr-util => ../go-addr-util

replace github.com/libp2p/go-buffer-pool => ../go-buffer-pool

replace github.com/libp2p/go-conn-security-multistream => ../go-conn-security-multistream

replace github.com/libp2p/go-eventbus => ../go-eventbus

replace github.com/libp2p/go-flow-metrics => ../go-flow-metrics

replace github.com/libp2p/go-libp2p => ../go-libp2p

replace github.com/libp2p/go-libp2p-autonat => ../go-libp2p-autonat

replace github.com/libp2p/go-libp2p-autonat-svc => ../go-libp2p-autonat-svc

replace github.com/libp2p/go-libp2p-blankhost => ../go-libp2p-blankhost

replace github.com/libp2p/go-libp2p-circuit => ../go-libp2p-circuit

replace github.com/libp2p/go-libp2p-connmgr => ../go-libp2p-connmgr

replace github.com/libp2p/go-libp2p-consensus => ../go-libp2p-consensus

replace github.com/libp2p/go-libp2p-daemon => ../go-libp2p-daemon

replace github.com/libp2p/go-libp2p-discovery => ../go-libp2p-discovery

replace github.com/libp2p/go-libp2p-examples => ../go-libp2p-examples

replace github.com/libp2p/go-libp2p-gorpc => ../go-libp2p-gorpc

replace github.com/libp2p/go-libp2p-introspection => ../go-libp2p-introspection

replace github.com/libp2p/go-libp2p-kad-dht => ../go-libp2p-kad-dht

replace github.com/libp2p/go-libp2p-kbucket => ../go-libp2p-kbucket

replace github.com/libp2p/go-libp2p-loggables => ../go-libp2p-loggables

replace github.com/libp2p/go-libp2p-mplex => ../go-libp2p-mplex

replace github.com/libp2p/go-libp2p-nat => ../go-libp2p-nat

replace github.com/libp2p/go-libp2p-netutil => ../go-libp2p-netutil

replace github.com/libp2p/go-libp2p-noise => ../go-libp2p-noise

replace github.com/libp2p/go-libp2p-peerstore => ../go-libp2p-peerstore

replace github.com/libp2p/go-libp2p-pnet => ../go-libp2p-pnet

replace github.com/libp2p/go-libp2p-pubsub => ../go-libp2p-pubsub

replace github.com/libp2p/go-libp2p-pubsub-router => ../go-libp2p-pubsub-router

replace github.com/libp2p/go-libp2p-pubsub-tracer => ../go-libp2p-pubsub-tracer

replace github.com/libp2p/go-libp2p-quic-transport => ../go-libp2p-quic-transport

replace github.com/libp2p/go-libp2p-raft => ../go-libp2p-raft

replace github.com/libp2p/go-libp2p-record => ../go-libp2p-record

replace github.com/libp2p/go-libp2p-routing-helpers => ../go-libp2p-routing-helpers

replace github.com/libp2p/go-libp2p-secio => ../go-libp2p-secio

replace github.com/libp2p/go-libp2p-swarm => ../go-libp2p-swarm

replace github.com/libp2p/go-libp2p-testing => ../go-libp2p-testing

replace github.com/libp2p/go-libp2p-tls => ../go-libp2p-tls

replace github.com/libp2p/go-libp2p-transport-upgrader => ../go-libp2p-transport-upgrader

replace github.com/libp2p/go-libp2p-webrtc-direct => ../go-libp2p-webrtc-direct

replace github.com/libp2p/go-libp2p-yamux => ../go-libp2p-yamux

replace github.com/libp2p/go-maddr-filter => ../go-maddr-filter

replace github.com/libp2p/go-mplex => ../go-mplex

replace github.com/libp2p/go-msgio => ../go-msgio

replace github.com/multiformats/go-multiaddr => ../go-multiaddr

replace github.com/multiformats/go-multiaddr-dns => ../go-multiaddr-dns

replace github.com/multiformats/go-multiaddr-fmt => ../go-multiaddr-fmt

replace github.com/multiformats/go-multiaddr-net => ../go-multiaddr-net

replace github.com/multiformats/go-multistream => ../go-multistream

replace github.com/libp2p/go-nat => ../go-nat

replace github.com/libp2p/go-reuseport => ../go-reuseport

replace github.com/libp2p/go-reuseport-transport => ../go-reuseport-transport

replace github.com/libp2p/go-sockaddr => ../go-sockaddr

replace github.com/libp2p/go-stream-muxer-multistream => ../go-stream-muxer-multistream

replace github.com/libp2p/go-tcp-transport => ../go-tcp-transport

replace github.com/libp2p/go-utp-transport => ../go-utp-transport

replace github.com/Jorropo/go-webrtc-aside-transport => ../go-webrtc-aside-transport

replace github.com/libp2p/go-ws-transport => ../go-ws-transport

replace github.com/libp2p/go-yamux => ../go-yamux
