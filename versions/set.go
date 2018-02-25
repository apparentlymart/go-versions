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
	// The special Unspecified version is excluded as soon as any sort of
	// constraint is applied, and so the only set it is a member of is
	// the special All set.
	if v == Unspecified {
		return s == All
	}

	return s.setI.Has(v)
}

var InitialDevelopment = OlderThan(MustParseVersion("1.0.0"))
