package helpers

// WithRace returns whether the binary was compiled
// with the race flag on.
func WithRace() bool {
	return withRace
}
