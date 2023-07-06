package versions

import (
	"encoding/json"
	"reflect"
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
			Union(
				Only(MustParseVersion("1.0.0")),
				Only(MustParseVersion("1.1.0")),
				Only(MustParseVersion("1.2.0")),
			),
			MustParseVersion("1.2.0"),
			true,
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
			All.Subtract(Only(MustParseVersion("0.9.0"))),
			MustParseVersion("0.9.1"),
			true,
		},
		{
			Union(
				All,
				Only(MustParseVersion("1.0.1")),
			).AllRequested(),
			MustParseVersion("1.0.1"),
			true,
		},
		{
			Union(
				All,
				Only(MustParseVersion("1.0.1")),
			).AllRequested(),
			MustParseVersion("1.0.2"),
			false,
		},
		{
			Intersection(
				All,
				Only(MustParseVersion("1.0.1")),
			).AllRequested(),
			MustParseVersion("1.0.1"),
			true,
		},
		{
			Intersection(
				All,
				Only(MustParseVersion("1.0.1")),
			).AllRequested(),
			MustParseVersion("1.0.2"),
			false,
		},
		{
			Intersection(
				AtLeast(MustParseVersion("2.0.0")),
				Only(MustParseVersion("1.0.1")),
			).AllRequested(),
			MustParseVersion("1.0.1"),
			false,
		},
		{
			Only(
				MustParseVersion("1.0.1"),
			).Subtract(
				AtLeast(MustParseVersion("1.0.0")),
			).AllRequested(),
			MustParseVersion("1.0.1"),
			false,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby(">= 1.0.0")),
			MustParseVersion("1.0.0"),
			true,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby(">= 1.0.0")),
			MustParseVersion("1.0.0-beta1"),
			false,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby(">= 1.0.0")),
			MustParseVersion("2.0.0-beta1"),
			false,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby("2.0.0-beta1")),
			MustParseVersion("2.0.0-beta1"),
			true,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby("!= 2.0.0-beta1, 2.0.0-beta1")),
			MustParseVersion("2.0.0-beta1"),
			false, // the constraint is contradictory, so includes nothing
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby(">= 1.0.0")),
			MustParseVersion("1.0.1"),
			true,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby(">= 1.0.0")).AllRequested(),
			MustParseVersion("0.0.1"),
			false,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby(">= 1.0.0")).AllRequested(),
			MustParseVersion("1.0.0-beta1"),
			false,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby(">= 1.0.0")).AllRequested(),
			MustParseVersion("2.0.0-beta1"),
			false,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby("2.0.0-beta1")).AllRequested(),
			MustParseVersion("2.0.0-beta1"),
			true,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby(">= 1.0.0")).WithoutUnrequestedPrereleases(),
			MustParseVersion("0.0.1"),
			false,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby(">= 1.0.0")).WithoutUnrequestedPrereleases(),
			MustParseVersion("1.0.0"),
			true,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby(">= 1.0.0")).WithoutUnrequestedPrereleases(),
			MustParseVersion("1.0.0-beta1"),
			false,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby(">= 1.0.0")).WithoutUnrequestedPrereleases(),
			MustParseVersion("2.0.0-beta1"),
			false,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby("2.0.0-beta1")).WithoutUnrequestedPrereleases(),
			MustParseVersion("2.0.0-beta1"),
			true,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby("~> 1.2.3")),
			MustParseVersion("1.2.3"),
			true,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby("~> 1.2.3")),
			MustParseVersion("1.2.5"),
			true,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby("~> 1.2.3")),
			MustParseVersion("1.3.0"),
			false,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby("~> 1.2")),
			MustParseVersion("1.2.3"),
			true,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby("~> 1.2")),
			MustParseVersion("1.2.5"),
			true,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby("~> 1.2")),
			MustParseVersion("1.3.0"),
			true,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby("~> 1.2")),
			MustParseVersion("2.0.0"),
			false,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby("~> 1")),
			MustParseVersion("1.2.3"),
			true,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby("~> 1")),
			MustParseVersion("1.2.5"),
			true,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby("~> 1")),
			MustParseVersion("1.3.0"),
			true,
		},
		{
			MustMakeSet(MeetingConstraintsStringRuby("~> 1")),
			MustParseVersion("2.0.0"),
			false,
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

func TestSetJSON(t *testing.T) {
	j := []byte(`"^1 || 2.0.0"`)
	var got Set
	err := json.Unmarshal(j, &got)
	if err != nil {
		t.Fatal(err)
	}

	want := Intersection(
		Released,
		Union(
			Intersection(
				AtLeast(MustParseVersion("1.0.0")),
				OlderThan(MustParseVersion("2.0.0")),
			),
			Only(MustParseVersion("2.0.0")),
		),
	)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("wrong result\ngot:  %#v\nwant :%#v", got, want)
	}
}
