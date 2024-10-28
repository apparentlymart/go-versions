package versions

import (
	"reflect"
	"testing"
)

func TestMeetingConstraintsCanon(t *testing.T) {
	tests := []struct {
		Input string
		Want  Set
	}{
		{
			`1.0.0`,
			Intersection(
				Released,
				Only(MustParseVersion(`1.0.0`)),
			),
		},
		{
			`=1.0.0`,
			Intersection(
				Released,
				Only(MustParseVersion(`1.0.0`)),
			),
		},
		{
			`1.0-beta.1`,
			// This result is sub-optimal since it mentions the pre-release
			// version twice, but it works. Perhaps later we'll try to
			// optimize this situation, but not too bothered for now.
			Union(
				Only(MustParseVersion(`1.0-beta.1`)),
				Intersection(
					Released,
					Only(MustParseVersion(`1.0-beta.1`)),
				),
			),
		},
		{
			`^1.0 || 2.0-beta.1 || 2.0-beta.2`,
			// This result is even less optimal, but again is functionally
			// correct.
			Union(
				Selection(
					MustParseVersion(`2.0-beta.1`),
					MustParseVersion(`2.0-beta.2`),
				),
				Intersection(
					Released,
					Union(
						Intersection(
							AtLeast(MustParseVersion("1.0.0")),
							OlderThan(MustParseVersion("2.0.0")),
						),
						Only(MustParseVersion(`2.0-beta.1`)),
						Only(MustParseVersion(`2.0-beta.2`)),
					),
				),
			),
		},
		{
			`!1.0.0`,
			Intersection(
				Released,
				All.Subtract(Only(MustParseVersion(`1.0.0`))),
			),
		},
		{
			`>1.0.0`,
			Intersection(
				Released,
				NewerThan(MustParseVersion(`1.0.0`)),
			),
		},
		{
			`>1.0`,
			Intersection(
				Released,
				NewerThan(MustParseVersion(`1.0.0`)),
			),
		},
		{
			`<1.0.0`,
			Intersection(
				Released,
				OlderThan(MustParseVersion(`1.0.0`)),
			),
		},
		{
			`>=1.0.0`,
			Intersection(
				Released,
				AtLeast(MustParseVersion(`1.0.0`)),
			),
		},
		{
			`<=1.0.0`,
			Intersection(
				Released,
				AtMost(MustParseVersion(`1.0.0`)),
			),
		},
		{
			`~1.2.3`,
			Intersection(
				Released,
				AtLeast(MustParseVersion(`1.2.3`)),
				OlderThan(MustParseVersion(`1.3.0`)),
			),
		},
		{
			`~1.2`,
			Intersection(
				Released,
				AtLeast(MustParseVersion(`1.2.0`)),
				OlderThan(MustParseVersion(`1.3.0`)),
			),
		},
		{
			`~1`,
			Intersection(
				Released,
				AtLeast(MustParseVersion(`1.0.0`)),
				OlderThan(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`^1.2.3`,
			Intersection(
				Released,
				AtLeast(MustParseVersion(`1.2.3`)),
				OlderThan(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`^0.2.3`,
			Intersection(
				Released,
				AtLeast(MustParseVersion(`0.2.3`)),
				OlderThan(MustParseVersion(`0.3.0`)),
			),
		},
		{
			`^1.2`,
			Intersection(
				Released,
				AtLeast(MustParseVersion(`1.2.0`)),
				OlderThan(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`^1`,
			Intersection(
				Released,
				AtLeast(MustParseVersion(`1.0.0`)),
				OlderThan(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`>=1 <2`,
			Intersection(
				Released,
				AtLeast(MustParseVersion(`1.0.0`)),
				OlderThan(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`1.*`,
			Intersection(
				Released,
				AtLeast(MustParseVersion(`1.0.0`)),
				OlderThan(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`1.2.*`,
			Intersection(
				Released,
				AtLeast(MustParseVersion(`1.2.0`)),
				OlderThan(MustParseVersion(`1.3.0`)),
			),
		},
		{
			`*`,
			Released,
		},
		{
			`*.*`,
			Released,
		},
		{
			`*.*.*`,
			Released,
		},
		{
			`1.0.0 2.0.0`,
			Intersection(
				Released,
				Only(MustParseVersion(`1.0.0`)),
				Only(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`>=1.0 || >=0.9 <0.10`,
			Intersection(
				Released,
				Union(
					AtLeast(MustParseVersion("1.0")),
					Intersection(
						AtLeast(MustParseVersion("0.9")),
						OlderThan(MustParseVersion("0.10")),
					),
				),
			),
		},
		{
			`1.0.0 1.0.0`, // redundant
			Intersection(
				Released,
				Only(MustParseVersion(`1.0.0`)), // the duplicates don't get optimized away (yet?)
				Only(MustParseVersion(`1.0.0`)), // probably not worth the effort but will test someday
			),
		},
		{
			`1.0.0 || 1.0.0`, // redundant
			Intersection(
				Released,
				Union(
					Only(MustParseVersion(`1.0.0`)), // the duplicates don't get optimized away (yet?)
					Only(MustParseVersion(`1.0.0`)), // probably not worth the effort but will test someday
				),
			),
		},
		{
			`1.0.0 !1.0.0`, // degenerate empty set
			Intersection(
				Released,
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
					t.Errorf("wrong result\ngot:  %+vwant: %+v", got, test.Want)
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
			Intersection(
				Released,
				Only(MustParseVersion(`1.0.0`)),
			),
		},
		{
			`= 1.0.0`,
			Intersection(
				Released,
				Only(MustParseVersion(`1.0.0`)),
			),
		},
		{
			`!= 1.0.0`,
			Intersection(
				Released,
				All.Subtract(Only(MustParseVersion(`1.0.0`))),
			),
		},
		{
			`> 1.0.0`,
			Intersection(
				Released,
				NewerThan(MustParseVersion(`1.0.0`)),
			),
		},
		{
			`> 1.0`,
			Intersection(
				Released,
				NewerThan(MustParseVersion(`1.0.0`)),
			),
		},
		{
			`< 1.0.0`,
			Intersection(
				Released,
				OlderThan(MustParseVersion(`1.0.0`)),
			),
		},
		{
			`>= 1.0.0`,
			Intersection(
				Released,
				AtLeast(MustParseVersion(`1.0.0`)),
			),
		},
		{
			`<= 1.0.0`,
			Intersection(
				Released,
				AtMost(MustParseVersion(`1.0.0`)),
			),
		},
		{
			`~> 1.2.3`,
			Intersection(
				Released,
				AtLeast(MustParseVersion(`1.2.3`)),
				OlderThan(MustParseVersion(`1.3.0`)),
			),
		},
		{
			`~> 1.2`,
			Intersection(
				Released,
				AtLeast(MustParseVersion(`1.2.0`)),
				OlderThan(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`~> 1`,
			Intersection(
				Released,
				AtLeast(MustParseVersion(`1.0.0`)),
				OlderThan(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`>= 1, < 2`,
			Intersection(
				Released,
				AtLeast(MustParseVersion(`1.0.0`)),
				OlderThan(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`1.0.0, 2.0.0`,
			Intersection(
				Released,
				Only(MustParseVersion(`1.0.0`)),
				Only(MustParseVersion(`2.0.0`)),
			),
		},
		{
			`1.0.0, 1.0.0`, // redundant
			Intersection(
				Released,
				Only(MustParseVersion(`1.0.0`)), // the duplicates don't get optimized away (yet?)
				Only(MustParseVersion(`1.0.0`)), // probably not worth the effort but will test someday
			),
		},
		{
			`1.0.0, != 1.0.0`, // degenerate empty set
			Intersection(
				Released,
				Only(MustParseVersion(`1.0.0`)),
				All.Subtract(Only(MustParseVersion(`1.0.0`))),
			),
		},
		{
			`1.0.0-beta1, != 1.0.0-beta1`, // degenerate empty set with prerelease versions
			Intersection(
				Released,
				Only(MustParseVersion(`1.0.0-beta1`)),
				All.Subtract(Only(MustParseVersion(`1.0.0-beta1`))),
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
					t.Errorf("wrong result\ngot:  %+vwant: %+v", got, test.Want)
				}
			}
		})
	}
}
