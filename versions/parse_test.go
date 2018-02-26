package versions

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestMeetingConstraintsCanon(t *testing.T) {
	tests := []struct {
		Input string
		Want  Set
	}{
		{
			`1.0.0`,
			Only(MustParseVersion(`1.0.0`)),
		},
		{
			`=1.0.0`,
			Only(MustParseVersion(`1.0.0`)),
		},
		{
			`1.0-beta.1`,
			Only(MustParseVersion(`1.0-beta.1`)),
		},
		{
			`!1.0.0`,
			All.Subtract(Only(MustParseVersion(`1.0.0`))),
		},
		{
			`>1.0.0`,
			NewerThan(MustParseVersion(`1.0.0`)),
		},
		{
			`>1.0`,
			NewerThan(MustParseVersion(`1.0.0`)),
		},
		{
			`<1.0.0`,
			OlderThan(MustParseVersion(`1.0.0`)),
		},
		{
			`>=1.0.0`,
			AtLeast(MustParseVersion(`1.0.0`)),
		},
		{
			`<=1.0.0`,
			AtMost(MustParseVersion(`1.0.0`)),
		},
		{
			`~1.2.3`,
			Intersection(
				AtLeast(MustParseVersion(`1.2.3`)),
				OlderThan(MustParseVersion(`1.3.0`)),
			),
		},
		{
			`~1.2`,
			Intersection(
				AtLeast(MustParseVersion(`1.2.0`)),
				OlderThan(MustParseVersion(`1.3.0`)),
			),
		},
		{
			`~1`,
			Intersection(
				AtLeast(MustParseVersion(`1.0.0`)),
				OlderThan(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`^1.2.3`,
			Intersection(
				AtLeast(MustParseVersion(`1.2.3`)),
				OlderThan(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`^0.2.3`,
			Intersection(
				AtLeast(MustParseVersion(`0.2.3`)),
				OlderThan(MustParseVersion(`0.3.0`)),
			),
		},
		{
			`^1.2`,
			Intersection(
				AtLeast(MustParseVersion(`1.2.0`)),
				OlderThan(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`^1`,
			Intersection(
				AtLeast(MustParseVersion(`1.0.0`)),
				OlderThan(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`>=1 <2`,
			Intersection(
				AtLeast(MustParseVersion(`1.0.0`)),
				OlderThan(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`1.*`,
			Intersection(
				AtLeast(MustParseVersion(`1.0.0`)),
				OlderThan(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`1.2.*`,
			Intersection(
				AtLeast(MustParseVersion(`1.2.0`)),
				OlderThan(MustParseVersion(`1.3.0`)),
			),
		},
		{
			`*`,
			All,
		},
		{
			`*.*`,
			All,
		},
		{
			`*.*.*`,
			All,
		},
		{
			`1.0.0 2.0.0`,
			Intersection(
				Only(MustParseVersion(`1.0.0`)),
				Only(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`>=1.0 || >=0.9 <0.10`,
			Union(
				AtLeast(MustParseVersion("1.0")),
				Intersection(
					AtLeast(MustParseVersion("0.9")),
					OlderThan(MustParseVersion("0.10")),
				),
			),
		},
		{
			`1.0.0 1.0.0`, // redundant
			Intersection(
				Only(MustParseVersion(`1.0.0`)), // the duplicates don't get optimized away (yet?)
				Only(MustParseVersion(`1.0.0`)), // probably not worth the effort but will test someday
			),
		},
		{
			`1.0.0 || 1.0.0`, // redundant
			Union(
				Only(MustParseVersion(`1.0.0`)), // the duplicates don't get optimized away (yet?)
				Only(MustParseVersion(`1.0.0`)), // probably not worth the effort but will test someday
			),
		},
		{
			`1.0.0 !1.0.0`, // degenerate empty set
			Intersection(
				Only(MustParseVersion(`1.0.0`)),
				All.Subtract(Only(MustParseVersion(`1.0.0`))),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			got, err := MeetingConstraintsString(test.Input)
			if err != nil {
				t.Fatal(err)
			}

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
func TestMeetingConstraintsRuby(t *testing.T) {
	tests := []struct {
		Input string
		Want  Set
	}{
		{
			`1.0.0`,
			Only(MustParseVersion(`1.0.0`)),
		},
		{
			`= 1.0.0`,
			Only(MustParseVersion(`1.0.0`)),
		},
		{
			`!= 1.0.0`,
			All.Subtract(Only(MustParseVersion(`1.0.0`))),
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
				Only(MustParseVersion(`1.0.0`)),
				Only(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`1.0.0, 1.0.0`, // redundant
			Intersection(
				Only(MustParseVersion(`1.0.0`)), // the duplicates don't get optimized away (yet?)
				Only(MustParseVersion(`1.0.0`)), // probably not worth the effort but will test someday
			),
		},
		{
			`1.0.0, != 1.0.0`, // degenerate empty set
			Intersection(
				Only(MustParseVersion(`1.0.0`)),
				All.Subtract(Only(MustParseVersion(`1.0.0`))),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			got, err := MeetingConstraintsStringRuby(test.Input)
			if err != nil {
				t.Fatal(err)
			}

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
