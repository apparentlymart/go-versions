package versions

// Set is a set of versions, usually created by parsing a constraint string.
type Set struct {
	setI
}

// setI is the private interface implemented by our various constraint
// operators.
type setI interface {
	Has(v Version) bool
	GoString() string
}

// Has returns true if the given version is a member of the receiving set.
func (s Set) Has(v Version) bool {
	// This wrapper is unnecessary because we embed the interface, but
	// we include it anyway so we can easily hang documentation from it.
	return s.setI.Has(v)
}
