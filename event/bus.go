package event

import "io"

// SubscriptionOpt represents a subscriber option. Use the options exposed by the implementation of choice.
type SubscriptionOpt = func(interface{}) error

// EmitterOpt represents an emitter option. Use the options exposed by the implementation of choice.
type EmitterOpt = func(interface{}) error

// CancelFunc closes a subscriber.
type CancelFunc = func()

// Emitter represents an actor that emits events onto the eventbus.
type Emitter interface {
	io.Closer

	// Emit emits an event onto the eventbus. If any channel subscribed to the topic is blocked,
	// calls to Emit will block.
	//
	// Calling this function with wrong event type will cause a panic.
	Emit(evt interface{})
}

type Subscription interface {
	Out() <-chan interface{}
	Close() error
}

// Bus is an interface for a type-based event delivery system.
type Bus interface {
	// Subscribe creates a new subscription.
	//
	// Failing to drain the channel may cause publishers to block. CancelFunc must return after
	// last send to the channel.
	//
	// Example:
	// sub, err := eventbus.Subscribe(new(EventType))
	// defer sub.Close()
	// for e := range sub.Out() {
	//   event := e.(EventType) // guaranteed safe
	//   [...]
	// }
	// TODO: update doc
	Subscribe(eventType interface{}, opts ...SubscriptionOpt) (Subscription, error)

	// Emitter creates a new event emitter.
	//
	// eventType accepts typed nil pointers, and uses the type information for wiring purposes.
	//
	// Example:
	// em, err := eventbus.Emitter(new(EventT))
	// defer em.Close() // MUST call this after being done with the emitter
	// em.Emit(EventT{})
	Emitter(eventType interface{}, opts ...EmitterOpt) (Emitter, error)
}
