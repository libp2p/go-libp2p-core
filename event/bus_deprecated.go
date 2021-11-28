package event

import (
	"github.com/libp2p/go-libp2p-core/eventbus"
)

// SubscriptionOpt represents a subscriber option. Use the options exposed by the implementation of choice.
// Deprecated: Use eventbus.SubscriptionOpt
type SubscriptionOpt = eventbus.SubscriptionOpt

// EmitterOpt represents an emitter option. Use the options exposed by the implementation of choice.
// Deprecated: Use eventbus.EmitterOpt
type EmitterOpt = eventbus.EmitterOpt

// CancelFunc closes a subscriber.
// Deprecated: Use eventbus.CancelFunc
type CancelFunc = eventbus.CancelFunc

// WildcardSubscription is the type to subscribe to to receive all events
// emitted in the eventbus.
// Deprecated: Use eventbus.WildcardSubscription
var WildcardSubscription = eventbus.WildcardSubscription

// Emitter represents an actor that emits events onto the eventbus.
// Deprecated: Use eventbus.Emitter
type Emitter = eventbus.Emitter

// Subscription represents a subscription to one or multiple event types.
// Deprecated: Use eventbus.Subscription
type Subscription = eventbus.Subscription

// Bus is an interface for a type-based event delivery system.
// Deprecated: Use eventbus.Bus
type Bus = eventbus.Bus
