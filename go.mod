module github.com/libp2p/go-libp2p-core

require (
	github.com/btcsuite/btcd v0.0.0-20190824003749-130ea5bddde3
	github.com/coreos/go-semver v0.3.0
	github.com/gogo/protobuf v1.3.1
	github.com/ipfs/go-cid v0.0.3
	github.com/jbenet/goprocess v0.1.3
	github.com/libp2p/go-flow-metrics v0.0.1
	github.com/libp2p/go-msgio v0.0.4
	github.com/libp2p/go-openssl v0.0.3
	github.com/minio/sha256-simd v0.1.1
	github.com/mr-tron/base58 v1.1.2
	github.com/multiformats/go-multiaddr v0.1.1
	github.com/multiformats/go-multihash v0.0.8
	github.com/smola/gocompat v0.2.0
	go.opencensus.io v0.22.1
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
)

//replace github.com/libp2p/go-flow-metrics v0.0.1 => github.com/kpp/go-flow-metrics v0.0.2-0.20191101005412-ce5ebda5e4a4 // Stebalien's version

replace github.com/libp2p/go-flow-metrics v0.0.1 => github.com/kpp/go-flow-metrics v0.0.2-0.20191031231915-edeb2d90f222 // kpp's version

go 1.12
