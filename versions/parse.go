package versions

import (
	"github.com/apparentlymart/go-versions/versions/constraints"
)

// ParseVersion attempts to parse the given string as a semantic version
// specification, and returns the result if successful.
//
// If the given string is not parseable then an error is returned that is
// suitable for display directly to a hypothetical end-user that provided this
// version string, as long as they can read English.
func ParseVersion(s string) (Version, error) {
	spec, err := constraints.ParseExactVersion(s)
	if err != nil {
		return Unspecified, err
	}
	return Version{
		Major:      spec.Major.Num,
		Minor:      spec.Minor.Num,
		Patch:      spec.Patch.Num,
		Prerelease: VersionExtra(spec.Prerelease),
		Metadata:   VersionExtra(spec.Metadata),
	}, nil
}
