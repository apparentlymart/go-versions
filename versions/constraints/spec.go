package constraints

type NumConstraint struct {
	Num           uint64
	Unconstrained bool
}

type VersionSpec struct {
	Major      NumConstraint
	Minor      NumConstraint
	Patch      NumConstraint
	Prerelease string
	Metadata   string
}
