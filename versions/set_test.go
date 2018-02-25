package versions

import (
	"testing"
)

// Make sure our set implementations all actually implement the interface
var _ setI = setBound{}
var _ setI = setExact{}
var _ setI = setExtreme(true)
var _ setI = setIntersection{}
var _ setI = setSubtract{}
var _ setI = setUnion{}
var _ setI = setReleased{}

func TestSetHas(t *testing.T) {
	tests := []struct {
		Set  Set
		Has  Version
		Want bool
	}{
		{
			All,
			Unspecified,
			true,
		},
		{
			None,
			Unspecified,
			false,
		},
		{
			All.Subtract(Only(MustParseVersion("1.0.0"))),
			Unspecified,
			false, // any sort of constraint removes the special Unspecified version
		},
		{
			InitialDevelopment,
			Unspecified,
			false,
		},
		{
			InitialDevelopment,
			MustParseVersion("0.0.2"),
			true,
		},
		{
			InitialDevelopment,
			MustParseVersion("1.0.0"),
			false,
		},
		{
			Released,
			MustParseVersion("1.0.0"),
			true,
		},
		{
			Released,
			MustParseVersion("1.0.0-beta1"),
			false,
		},
		{
			Prerelease,
			MustParseVersion("1.0.0"),
			false,
		},
		{
			Prerelease,
			MustParseVersion("1.0.0-beta1"),
			true,
		},
		{
			Union(
				Only(MustParseVersion("1.0.0")),
				Only(MustParseVersion("1.1.0")),
			),
			MustParseVersion("1.0.0"),
			true,
		},
		{
			Union(
				Only(MustParseVersion("1.0.0")),
				Only(MustParseVersion("1.1.0")),
			),
			MustParseVersion("1.1.0"),
			true,
		},
		{
			Union(
				Only(MustParseVersion("1.0.0")),
				Only(MustParseVersion("1.1.0")),
			),
			MustParseVersion("1.2.0"),
			false,
		},
		{
			Intersection(
				AtLeast(MustParseVersion("1.0.0")),
				OlderThan(MustParseVersion("2.0.0")),
			),
			MustParseVersion("0.0.2"),
			false,
		},
		{
			Intersection(
				AtLeast(MustParseVersion("1.0.0")),
				OlderThan(MustParseVersion("2.0.0")),
			),
			MustParseVersion("1.0.0"),
			true,
		},
		{
			Intersection(
				AtLeast(MustParseVersion("1.0.0")),
				OlderThan(MustParseVersion("2.0.0")),
			),
			MustParseVersion("1.2.3"),
			true,
		},
		{
			Intersection(
				AtLeast(MustParseVersion("1.0.0")),
				OlderThan(MustParseVersion("2.0.0")),
			),
			MustParseVersion("2.0.0"),
			false,
		},
		{
			Intersection(
				AtLeast(MustParseVersion("1.0.0")),
				OlderThan(MustParseVersion("2.0.0")),
			),
			MustParseVersion("2.0.1"),
			false,
		},
		{
			All.Subtract(Only(MustParseVersion("0.9.0"))),
			MustParseVersion("0.9.0"),
			false,
		},
		{
			All.Subtract(Only(MustParseVersion("0.9.0"))),
			MustParseVersion("0.9.1"),
			true,
		},
		{
			MustMakeSet(MeetingRubyStyleConstraints(">= 1.0.0")),
			MustParseVersion("1.0.0"),
			true,
		},
		{
			MustMakeSet(MeetingRubyStyleConstraints(">= 1.0.0")),
			MustParseVersion("1.0.0-beta1"),
			false,
		},
		{
			MustMakeSet(MeetingRubyStyleConstraints(">= 1.0.0")),
			MustParseVersion("2.0.0-beta1"),
			true,
		},
		{
			MustMakeSet(MeetingRubyStyleConstraints(">= 1.0.0")),
			MustParseVersion("1.0.1"),
			true,
		},
		{
			MustMakeSet(MeetingRubyStyleConstraints(">= 1.0.0")),
			MustParseVersion("0.0.1"),
			false,
		},
		{
			MustMakeSet(RequestedByRubyStyleConstraints(">= 1.0.0")),
			MustParseVersion("1.0.0-beta1"),
			false,
		},
		{
			MustMakeSet(RequestedByRubyStyleConstraints(">= 1.0.0")),
			MustParseVersion("2.0.0-beta1"),
			false,
		},
		{
			MustMakeSet(RequestedByRubyStyleConstraints("2.0.0-beta1")),
			MustParseVersion("2.0.0-beta1"),
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.Set.GoString(), func(t *testing.T) {
			got := test.Set.Has(test.Has)

			if got != test.Want {
				t.Errorf(
					"wrong result\nset:     %#v\nversion: %#v\ngot:     %#v\nwant:    %#v",
					test.Set,
					test.Has,
					got, test.Want,
				)
			}
		})
	}
}
