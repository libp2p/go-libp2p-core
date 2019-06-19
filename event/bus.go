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

// Bus is an interface to type-based event delivery system
type Bus interface {
	// Subscribe creates new subscription. Failing to drain the channel will cause
	// publishers to get blocked. CancelFunc is guaranteed to return after last send
	// to the channel
	//
	// Example:
	// ch := make(chan EventT, 10)
	// defer close(ch)
	// cancel, err := eventbus.Subscribe(ch)
	// defer cancel()
	Subscribe(typedChan interface{}, opts ...SubscriptionOpt) (CancelFunc, error)

	// Emitter creates new emitter
	//
	// eventType accepts typed nil pointers, and uses the type information to
	// select output type
	//
	// Example:
	// em, err := eventbus.Emitter(new(EventT))
	// defer em.Close() // MUST call this after being done with the emitter
	//
	// em.Emit(EventT{})
	Emitter(eventType interface{}, opts ...EmitterOpt) (Emitter, error)
}
