package constraints

import (
	"testing"

	"github.com/go-test/deep"
)

func TestParseRubyStyle(t *testing.T) {
	tests := []struct {
		Input   string
		Want    SelectionSpec
		WantErr string
	}{
		{
			"",
			SelectionSpec{},
			"empty specification",
		},
		{
			"1",
			SelectionSpec{
				Operator: OpEqual,
				Boundary: VersionSpec{
					Major: NumConstraint{Num: 1},
					Minor: NumConstraint{Unconstrained: true},
					Patch: NumConstraint{Unconstrained: true},
				},
			},
			"",
		},
		{
			"1.1",
			SelectionSpec{
				Operator: OpEqual,
				Boundary: VersionSpec{
					Major: NumConstraint{Num: 1},
					Minor: NumConstraint{Num: 1},
					Patch: NumConstraint{Unconstrained: true},
				},
			},
			"",
		},
		{
			"1.1.1",
			SelectionSpec{
				Operator: OpEqual,
				Boundary: VersionSpec{
					Major: NumConstraint{Num: 1},
					Minor: NumConstraint{Num: 1},
					Patch: NumConstraint{Num: 1},
				},
			},
			"",
		},
		{
			"1.0.0.0",
			SelectionSpec{},
			"too many numbered portions; only three are allowed (major, minor, patch)",
		},
		{
			"v1.0.0",
			SelectionSpec{},
			`a "v" prefix should not be used`,
		},
		{
			"1.0.0-beta2",
			SelectionSpec{
				Operator: OpEqual,
				Boundary: VersionSpec{
					Major:      NumConstraint{Num: 1},
					Minor:      NumConstraint{Num: 0},
					Patch:      NumConstraint{Num: 0},
					Prerelease: "beta2",
				},
			},
			"",
		},
		{
			"1.0-beta2",
			SelectionSpec{
				Operator: OpEqual,
				Boundary: VersionSpec{
					Major:      NumConstraint{Num: 1},
					Minor:      NumConstraint{Num: 0},
					Patch:      NumConstraint{Num: 0}, // implied by the prerelease tag to ensure constraint consistency
					Prerelease: "beta2",
				},
			},
			"",
		},
		{
			"1.0.0-beta.2",
			SelectionSpec{
				Operator: OpEqual,
				Boundary: VersionSpec{
					Major:      NumConstraint{Num: 1},
					Minor:      NumConstraint{Num: 0},
					Patch:      NumConstraint{Num: 0},
					Prerelease: "beta.2",
				},
			},
			"",
		},
		{
			"1.0.0+foo",
			SelectionSpec{
				Operator: OpEqual,
				Boundary: VersionSpec{
					Major:    NumConstraint{Num: 1},
					Minor:    NumConstraint{Num: 0},
					Patch:    NumConstraint{Num: 0},
					Metadata: "foo",
				},
			},
			"",
		},
		{
			"1.0.0+foo.bar",
			SelectionSpec{
				Operator: OpEqual,
				Boundary: VersionSpec{
					Major:    NumConstraint{Num: 1},
					Minor:    NumConstraint{Num: 0},
					Patch:    NumConstraint{Num: 0},
					Metadata: "foo.bar",
				},
			},
			"",
		},
		{
			"1.0.0-beta1+foo.bar",
			SelectionSpec{
				Operator: OpEqual,
				Boundary: VersionSpec{
					Major:      NumConstraint{Num: 1},
					Minor:      NumConstraint{Num: 0},
					Patch:      NumConstraint{Num: 0},
					Prerelease: "beta1",
					Metadata:   "foo.bar",
				},
			},
			"",
		},
		{
			"> 1.1.1",
			SelectionSpec{
				Operator: OpGreaterThan,
				Boundary: VersionSpec{
					Major: NumConstraint{Num: 1},
					Minor: NumConstraint{Num: 1},
					Patch: NumConstraint{Num: 1},
				},
			},
			``,
		},
		{
			">= 1.1.1",
			SelectionSpec{
				Operator: OpGreaterThanOrEqual,
				Boundary: VersionSpec{
					Major: NumConstraint{Num: 1},
					Minor: NumConstraint{Num: 1},
					Patch: NumConstraint{Num: 1},
				},
			},
			``,
		},
		{
			"=> 1.1.1",
			SelectionSpec{},
			`invalid constraint operator "=>"; did you mean ">="?`,
		},
		{
			"< 1.1.1",
			SelectionSpec{
				Operator: OpLessThan,
				Boundary: VersionSpec{
					Major: NumConstraint{Num: 1},
					Minor: NumConstraint{Num: 1},
					Patch: NumConstraint{Num: 1},
				},
			},
			``,
		},
		{
			"<= 1.1.1",
			SelectionSpec{
				Operator: OpLessThanOrEqual,
				Boundary: VersionSpec{
					Major: NumConstraint{Num: 1},
					Minor: NumConstraint{Num: 1},
					Patch: NumConstraint{Num: 1},
				},
			},
			``,
		},
		{
			"=< 1.1.1",
			SelectionSpec{},
			`invalid constraint operator "=<"; did you mean "<="?`,
		},
		{
			"~> 1.1.1",
			SelectionSpec{
				Operator: OpGreaterThanOrEqualPatchOnly,
				Boundary: VersionSpec{
					Major: NumConstraint{Num: 1},
					Minor: NumConstraint{Num: 1},
					Patch: NumConstraint{Num: 1},
				},
			},
			``,
		},
		{
			"~> 1.1",
			SelectionSpec{
				Operator: OpGreaterThanOrEqualMinorOnly,
				Boundary: VersionSpec{
					Major: NumConstraint{Num: 1},
					Minor: NumConstraint{Num: 1},
					Patch: NumConstraint{Unconstrained: true},
				},
			},
			``,
		},
		{
			"~> 1",
			SelectionSpec{
				Operator: OpGreaterThanOrEqualMinorOnly,
				Boundary: VersionSpec{
					Major: NumConstraint{Num: 1},
					Minor: NumConstraint{Unconstrained: true},
					Patch: NumConstraint{Unconstrained: true},
				},
			},
			``,
		},
		{
			"= 1.1.1",
			SelectionSpec{
				Operator: OpEqual,
				Boundary: VersionSpec{
					Major: NumConstraint{Num: 1},
					Minor: NumConstraint{Num: 1},
					Patch: NumConstraint{Num: 1},
				},
			},
			``,
		},
		{
			"!= 1.1.1",
			SelectionSpec{
				Operator: OpNotEqual,
				Boundary: VersionSpec{
					Major: NumConstraint{Num: 1},
					Minor: NumConstraint{Num: 1},
					Patch: NumConstraint{Num: 1},
				},
			},
			``,
		},
		{
			"=1.1.1",
			SelectionSpec{},
			`a space separator is required after the operator "="`,
		},
		{
			"=  1.1.1",
			SelectionSpec{},
			`only one space is expected after the operator "="`,
		},
		{
			"garbage",
			SelectionSpec{},
			`invalid characters "garbage"`,
		},
		{
			"& 1.1.0",
			SelectionSpec{},
			`invalid constraint operator "&"`,
		},
		{
			"1.*.*",
			SelectionSpec{},
			`can't use wildcard for minor number; omit segments that should be unconstrained`,
		},
		{
			"1.0.x",
			SelectionSpec{},
			`can't use wildcard for patch number; omit segments that should be unconstrained`,
		},
		{
			"1.0 || 2.0",
			SelectionSpec{},
			`only one constraint may be specified`,
		},
		{
			"1.0.0, 2.0.0",
			SelectionSpec{},
			`only one constraint may be specified`,
		},
		{
			"1.0.0 - 2.0.0",
			SelectionSpec{},
			`range constraints are not supported`,
		},
	}

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			got, err := ParseRubyStyle(test.Input)
			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}
			if gotErr != test.WantErr {
				t.Errorf("wrong error\ngot:  %s\nwant: %s", gotErr, test.WantErr)
				return
			}
			if err != nil {
				return
			}

			for _, problem := range deep.Equal(got, test.Want) {
				t.Error(problem)
			}
		})
	}
}

func TestParseRubyStyleMulti(t *testing.T) {
	tests := []struct {
		Input   string
		Want    IntersectionSpec
		WantErr string
	}{
		{
			"",
			nil,
			"",
		},
		{
			"1.1.1",
			IntersectionSpec{
				SelectionSpec{
					Operator: OpEqual,
					Boundary: VersionSpec{
						Major: NumConstraint{Num: 1},
						Minor: NumConstraint{Num: 1},
						Patch: NumConstraint{Num: 1},
					},
				},
			},
			"",
		},
		{
			">= 1.0, < 2",
			IntersectionSpec{
				SelectionSpec{
					Operator: OpGreaterThanOrEqual,
					Boundary: VersionSpec{
						Major: NumConstraint{Num: 1},
						Minor: NumConstraint{Num: 0},
						Patch: NumConstraint{Unconstrained: true},
					},
				},
				SelectionSpec{
					Operator: OpLessThan,
					Boundary: VersionSpec{
						Major: NumConstraint{Num: 2},
						Minor: NumConstraint{Unconstrained: true},
						Patch: NumConstraint{Unconstrained: true},
					},
				},
			},
			"",
		},
		{
			">= 1.0 < 2",
			nil,
			`missing comma after ">= 1.0"`,
		},
	}

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			got, err := ParseRubyStyleMulti(test.Input)
			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}
			if gotErr != test.WantErr {
				t.Errorf("wrong error\ngot:  %s\nwant: %s", gotErr, test.WantErr)
				return
			}
			if err != nil {
				return
			}

			for _, problem := range deep.Equal(got, test.Want) {
				t.Error(problem)
			}
		})
	}
}
