package constraints

import (
	"bytes"
	"fmt"
	"strconv"
)

type VersionSpec struct {
	Major      NumConstraint
	Minor      NumConstraint
	Patch      NumConstraint
	Prerelease string
	Metadata   string
}

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
