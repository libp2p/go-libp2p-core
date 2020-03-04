package peerstore

// WriteOpt is an option to be used in peerstore write operations.
type WriteOpt func(*writeOpts) error

// ReadOpt is an option to be used in peerstore read/query operations.
type ReadOpt func(*readOpts) error

type writeOpts struct {
	labels []Label
}

type readOpts struct {
	includeLabels []Label
	excludeLabels []Label
}

// Labels is a write option to set labels on addresses or peers during a write
// operation.
func Labels(labels ...Label) WriteOpt {
	return func(wo *writeOpts) error {
		wo.labels = labels
		return nil
	}
}

// IncludeLabels is a read option that restricts the results of a read/query
// operation to include ONLY addresses or peers with the listed labels.
func IncludeLabels(labels ...Label) ReadOpt {
	return func(ro *readOpts) error {
		ro.includeLabels = labels
		return nil
	}
}

// ExcludeLabels is a read option that restricts the results of a read/query
// operation to include all hits BUT the ones with the listed labels.
func ExcludeLabels(labels ...Label) ReadOpt {
	return func(ro *readOpts) error {
		ro.excludeLabels = labels
		return nil
	}
}
