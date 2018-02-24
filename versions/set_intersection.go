package versions

type setIntersection []setI

func (s setIntersection) Has(v Version) bool {
	if len(s) == 0 {
		// Weird to have an intersection with no elements, but we'll
		// allow it and return something sensible.
		return false
	}
	for _, ss := range s {
		if !ss.Has(v) {
			return false
		}
	}
	return true
}

// Intersection returns a new set that contains all of the versions that
// the receiver and the given sets have in common.
//
// The result is a finite set if the receiver or any of the given sets are
// finite.
func (s Set) Intersection(others ...Set) Set {
	r := make(setIntersection, 1, len(others)+1)
	r[0] = s.setI
	for _, ss := range others {
		if ss.setI == All {
			continue
		}
		if su, ok := ss.setI.(setIntersection); ok {
			r = append(r, su...)
		} else {
			r = append(r, ss.setI)
		}
	}
	return Set{setI: r}
}

var _ setFinite = setIntersection{}

func (s setIntersection) isFinite() bool {
	// intersection is finite if any of its members are, or if it is empty
	if len(s) == 0 {
		return true
	}
	for _, ss := range s {
		if isFinite(ss) {
			return true
		}
	}
	return false
}

func (s setIntersection) listVersions() List {
	var ret List
	for _, ss := range s {
		if isFinite(ss) {
			ret = append(ret, ss.(setFinite).listVersions()...)
		}
	}
	ret.Filter(Set{setI: s})
	return ret
}
