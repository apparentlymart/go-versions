package versions

type setExact map[Version]struct{}

func (s setExact) Has(v Version) bool {
	_, has := s[v]
	return has
}

// ExactVersion returns a version set containing only the given version.
//
// This function is guaranteed to produce a finite set.
func Exactly(v Version) Set {
	return Set{
		setI: setExact{v: struct{}{}},
	}
}

// Selection returns a version set containing only the versions given
// as arguments.
//
// This function is guaranteed to produce a finite set.
func Selection(vs ...Version) Set {
	ret := make(setExact)
	for _, v := range vs {
		ret[v] = struct{}{}
	}
	return Set{setI: ret}
}

// Exactly returns true if and only if the receiving set is finite and
// contains only a single version that is the same as the version given.
func (s Set) Exactly(v Version) bool {
	if !s.IsFinite() {
		return false
	}
	l := s.List()
	if len(l) != 1 {
		return false
	}
	return v.Same(l[0])
}

// Selects returns true if and only if the receiving set is finite and
// one of its selected versions is the given version. This is a weaker
// version of Exactly that allows the given version to be one of many
// exact versions.
func (s Set) Selects(v Version) bool {
	if !s.IsFinite() {
		return false
	}
	for _, pv := range s.List() {
		if v.Same(pv) {
			return true
		}
	}
	return false
}

var _ setFinite = setExact(nil)

func (s setExact) isFinite() bool {
	return true
}

func (s setExact) listVersions() List {
	if len(s) == 0 {
		return nil
	}
	ret := make(List, 0, len(s))
	for v := range s {
		ret = append(ret, v)
	}
	return ret
}
