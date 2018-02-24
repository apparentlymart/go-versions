package versions

type setUnion []setI

func (s setUnion) Has(v Version) bool {
	for _, ss := range s {
		if ss.Has(v) {
			return true
		}
	}
	return false
}

// Union returns a new set that contains all of the versions from the
// receiver and all of the versions from each of the other given sets.
//
// The result is finite only if the receiver and all of the other given sets
// are finite.
func (s Set) Union(others ...Set) Set {
	r := make(setUnion, 1, len(others)+1)
	r[0] = s.setI
	for _, ss := range others {
		if ss.setI == None {
			continue
		}
		if su, ok := ss.setI.(setUnion); ok {
			r = append(r, su...)
		} else {
			r = append(r, ss.setI)
		}
	}
	return Set{setI: r}
}

var _ setFinite = setUnion{}

func (s setUnion) isFinite() bool {
	// union is finite only if all of its members are finite
	for _, ss := range s {
		if !isFinite(ss) {
			return false
		}
	}
	return true
}

func (s setUnion) listVersions() List {
	var ret List
	for _, ss := range s {
		ret = append(ret, ss.(setFinite).listVersions()...)
	}
	return ret
}
