package discovery

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
)

// Advertiser is an interface for advertising services
type Advertiser interface {
	// Advertise advertises a service
	Advertise(ctx context.Context, ns string, opts ...DiscoveryOpt) (time.Duration, error)
}

// Discoverer is an interface for peer discovery
type Discoverer interface {
	// FindPeers discovers peers providing a service
	FindPeers(ctx context.Context, ns string, opts ...DiscoveryOpt) (<-chan peer.Info, error)
}

// Discovery is an interface that combines service advertisement and peer discovery
type Discovery interface {
	Advertiser
	Discoverer
}

// DiscoveryOpt is a single discovery option.
type DiscoveryOpt func(opts *DiscoveryOpts) error

// DiscoveryOpts is a set of discovery options.
type DiscoveryOpts struct {
	Ttl   time.Duration
	Limit int

	// Other (implementation-specific) options
	Other map[interface{}]interface{}
}

// Apply applies the given options to this DiscoveryOpts
func (opts *DiscoveryOpts) Apply(options ...DiscoveryOpt) error {
	for _, o := range options {
		if err := o(opts); err != nil {
			return err
		}
	}
	return nil
}
