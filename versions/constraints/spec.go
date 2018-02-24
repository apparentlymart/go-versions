package constraints

import (
	"bytes"
	"fmt"
	"strconv"
)

// Spec is an interface type that UnionSpec, IntersectionSpec, SelectionSpec,
// and VersionSpec all belong to.
//
// It's provided to allow generic code to be written that accepts and operates
// on all specs, but such code must still handle each type separately using
// e.g. a type switch. This is a closed type that will not have any new
// implementations added in future.
type Spec interface {
	isSpec()
}

// UnionSpec represents an "or" operation on nested version constraints.
//
// This is not directly representable in all of our supported constraint
// syntaxes.
type UnionSpec []IntersectionSpec

func (s UnionSpec) isSpec() {}

// IntersectionSpec represents an "and" operation on nested version constraints.
type IntersectionSpec []SelectionSpec

func (s IntersectionSpec) isSpec() {}

// SelectionSpec represents applying a single operator to a particular
// "boundary" version.
type SelectionSpec struct {
	Boundary VersionSpec
	Operator SelectionOp
}

func (s SelectionSpec) isSpec() {}

// VersionSpec represents the boundary within a SelectionSpec.
type VersionSpec struct {
	Major      NumConstraint
	Minor      NumConstraint
	Patch      NumConstraint
	Prerelease string
	Metadata   string
}

func (s VersionSpec) isSpec() {}

func (s VersionSpec) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s.%s.%s", s.Major, s.Minor, s.Patch)
	if s.Prerelease != "" {
		fmt.Fprintf(&buf, "-%s", s.Prerelease)
	}
	if s.Metadata != "" {
		fmt.Fprintf(&buf, "+%s", s.Metadata)
	}
	return buf.String()
}

type SelectionOp rune

const OpGreaterThan = '>'
const OpLessThan = '<'
const OpGreaterThanOrEqual = '≥'
const OpGreaterThanOrEqualPatchOnly = '~'
const OpGreaterThanOrEqualMinorOnly = '^'
const OpLessThanOrEqual = '≤'
const OpEqual = '='
const OpNotEqual = '≠'
const OpWildcards = '*'

type NumConstraint struct {
	Num           uint64
	Unconstrained bool
}

func (c NumConstraint) String() string {
	if c.Unconstrained {
		return "*"
	} else {
		return strconv.FormatUint(c.Num, 10)
	}
}
