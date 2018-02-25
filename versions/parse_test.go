package versions

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/apparentlymart/go-versions/versions/constraints"
)

func TestMeetingConstraints(t *testing.T) {
	// We're using the "ruby-like" constraint syntax here just because
	// it's nice and terse for our test table. This logic should work
	// regardless of which parser is used.

	tests := []struct {
		Input string
		Want  Set
	}{
		{
			`1.0.0`,
			Exactly(MustParseVersion(`1.0.0`)),
		},
		{
			`= 1.0.0`,
			Exactly(MustParseVersion(`1.0.0`)),
		},
		{
			`!= 1.0.0`,
			All.Subtract(Exactly(MustParseVersion(`1.0.0`))),
		},
		{
			`> 1.0.0`,
			NewerThan(MustParseVersion(`1.0.0`)),
		},
		{
			`> 1.0`,
			NewerThan(MustParseVersion(`1.0.0`)),
		},
		{
			`< 1.0.0`,
			OlderThan(MustParseVersion(`1.0.0`)),
		},
		{
			`>= 1.0.0`,
			AtLeast(MustParseVersion(`1.0.0`)),
		},
		{
			`<= 1.0.0`,
			AtMost(MustParseVersion(`1.0.0`)),
		},
		{
			`~> 1.2.3`,
			Intersection(
				AtLeast(MustParseVersion(`1.2.3`)),
				OlderThan(MustParseVersion(`1.3.0`)),
			),
		},
		{
			`~> 1.2`,
			Intersection(
				AtLeast(MustParseVersion(`1.2.0`)),
				OlderThan(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`~> 1`,
			Intersection(
				AtLeast(MustParseVersion(`1.0.0`)),
				OlderThan(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`>= 1, < 2`,
			Intersection(
				AtLeast(MustParseVersion(`1.0.0`)),
				OlderThan(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`1.0.0, 2.0.0`,
			Intersection(
				Exactly(MustParseVersion(`1.0.0`)),
				Exactly(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`1.0.0, 1.0.0`, // redundant
			Intersection(
				Exactly(MustParseVersion(`1.0.0`)), // the duplicates don't get optimized away (yet?)
				Exactly(MustParseVersion(`1.0.0`)), // probably not worth the effort but will test someday
			),
		},
		{
			`1.0.0, != 1.0.0`, // degenerate empty set
			Intersection(
				Exactly(MustParseVersion(`1.0.0`)),
				All.Subtract(Exactly(MustParseVersion(`1.0.0`))),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			spec, err := constraints.ParseRubyStyleMulti(test.Input)
			if err != nil {
				t.Fatal(err)
			}

			got := MeetingConstraints(spec)
			if !reflect.DeepEqual(got, test.Want) {
				gotStr := got.GoString()
				wantStr := test.Want.GoString()
				if gotStr != wantStr {
					t.Errorf("wrong result\ngot:  %s\nwant: %s", gotStr, wantStr)
				} else {
					// Sometimes our GoString implementations hide differences that
					// DeepEqual thinks are significant.
					t.Errorf("wrong result\ngot:  %swant: %s", spew.Sdump(got), spew.Sdump(test.Want))
				}
			}
		})
	}
}
