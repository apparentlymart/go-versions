package versions

type setSubtract struct {
	from setI
	sub  setI
}

func (s setSubtract) Has(v Version) bool {
	return s.from.Has(v) && !s.sub.Has(v)
}

// Subtract returns a new set that has all of the versions from the receiver
// except for any versions in the other given set.
//
// If the receiver is finite then the returned set is also finite.
func (s Set) Subtract(other Set) Set {
	if other.setI == None || s.setI == None {
		return s
	}
	if other.setI == All {
		return None
	}
	return Set{
		setI: setSubtract{
			from: s,
			sub:  other,
		},
	}
}

var _ setFinite = setSubtract{}

func (s setSubtract) isFinite() bool {
	// subtract is finite if its "from" is finite
	return isFinite(s.from)
}

func (s setSubtract) listVersions() List {
	ret := s.from.(setFinite).listVersions()
	ret = ret.Filter(Set{setI: s.sub})
	return ret
}
