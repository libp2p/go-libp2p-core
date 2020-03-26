package event

// EvtLocalHostInitialized is emitted when the local Host has been fully initialized.
// Once this event is emitted by the Host, subscribers are guaranteed that the Host
// has completely finished initializing and will not instantiate any new components.
type EvtLocalHostInitialized struct{}
