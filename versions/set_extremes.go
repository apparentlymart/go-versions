package versions

// All is an infinite set containing all possible versions.
var All Set

// None is a finite set containing no versions.
var None Set

type setExtreme bool

func (s setExtreme) Has(v Version) bool {
	return bool(s)
}

var _ setFinite = setExtreme(false)

func (s setExtreme) isFinite() bool {
	// Only None is finite
	return !bool(s)
}

func (s setExtreme) listVersions() List {
	return nil
}

func init() {
	All = Set{
		setI: setExtreme(true),
	}
	None = Set{
		setI: setExtreme(false),
	}
}
