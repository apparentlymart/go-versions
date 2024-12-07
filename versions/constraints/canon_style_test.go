package constraints

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-test/deep"
)

func TestParse(t *testing.T) {
	tests := []struct {
		Input   string
		Want    UnionSpec
		WantErr string
	}{
		{
			"",
			nil,
			"empty specification",
		},
		{
			"1",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpEqual,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Num: 0},
							Patch: NumConstraint{Num: 0},
						},
					},
				},
			},
			"",
		},
		{
			"1.1",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpEqual,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Num: 1},
							Patch: NumConstraint{Num: 0},
						},
					},
				},
			},
			"",
		},
		{
			"1.1.1",
			UnionSpec{
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
			},
			"",
		},
		{
			"1.0.0.0",
			nil,
			"too many numbered portions; only three are allowed (major, minor, patch)",
		},
		{
			"v1.0.0",
			nil,
			`a "v" prefix should not be used when specifying versions`,
		},
		{
			"1.0.0-beta2",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpEqual,
						Boundary: VersionSpec{
							Major:      NumConstraint{Num: 1},
							Minor:      NumConstraint{Num: 0},
							Patch:      NumConstraint{Num: 0},
							Prerelease: "beta2",
						},
					},
				},
			},
			"",
		},
		{
			"1.0-beta2",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpEqual,
						Boundary: VersionSpec{
							Major:      NumConstraint{Num: 1},
							Minor:      NumConstraint{Num: 0},
							Patch:      NumConstraint{Num: 0}, // implied by the prerelease tag to ensure constraint consistency
							Prerelease: "beta2",
						},
					},
				},
			},
			"",
		},
		{
			"1.0.0-beta.2",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpEqual,
						Boundary: VersionSpec{
							Major:      NumConstraint{Num: 1},
							Minor:      NumConstraint{Num: 0},
							Patch:      NumConstraint{Num: 0},
							Prerelease: "beta.2",
						},
					},
				},
			},
			"",
		},
		{
			"1.0.0+foo",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpEqual,
						Boundary: VersionSpec{
							Major:    NumConstraint{Num: 1},
							Minor:    NumConstraint{Num: 0},
							Patch:    NumConstraint{Num: 0},
							Metadata: "foo",
						},
					},
				},
			},
			"",
		},
		{
			"1.0.0+foo.bar",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpEqual,
						Boundary: VersionSpec{
							Major:    NumConstraint{Num: 1},
							Minor:    NumConstraint{Num: 0},
							Patch:    NumConstraint{Num: 0},
							Metadata: "foo.bar",
						},
					},
				},
			},
			"",
		},
		{
			"1.0.0-beta1+foo.bar",
			UnionSpec{
				IntersectionSpec{
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
				},
			},
			"",
		},
		{
			">1.1.1",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpGreaterThan,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Num: 1},
							Patch: NumConstraint{Num: 1},
						},
					},
				},
			},
			``,
		},
		{
			">2.*.*",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpGreaterThanOrEqual,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 3},
							Minor: NumConstraint{Num: 0},
							Patch: NumConstraint{Num: 0},
						},
					},
				},
			},
			``,
		},
		{
			">=1.1.1",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpGreaterThanOrEqual,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Num: 1},
							Patch: NumConstraint{Num: 1},
						},
					},
				},
			},
			``,
		},
		{
			">=2.*.*",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpGreaterThanOrEqual,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 2},
							Minor: NumConstraint{Num: 0},
							Patch: NumConstraint{Num: 0},
						},
					},
				},
			},
			``,
		},
		{
			"=>1.1.1",
			nil,
			`invalid constraint operator "=>"; did you mean ">="?`,
		},
		{
			"<1.1.1",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpLessThan,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Num: 1},
							Patch: NumConstraint{Num: 1},
						},
					},
				},
			},
			``,
		},
		{
			"<2.*.*",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpLessThan,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 2},
							Minor: NumConstraint{Num: 0},
							Patch: NumConstraint{Num: 0},
						},
					},
				},
			},
			``,
		},
		{
			"<=1.1.1",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpLessThanOrEqual,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Num: 1},
							Patch: NumConstraint{Num: 1},
						},
					},
				},
			},
			``,
		},
		{
			"<=2.*.*",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpLessThan,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 3},
							Minor: NumConstraint{Num: 0},
							Patch: NumConstraint{Num: 0},
						},
					},
				},
			},
			``,
		},
		{
			"=<1.1.1",
			nil,
			`invalid constraint operator "=<"; did you mean "<="?`,
		},
		{
			"~1.1.1",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpGreaterThanOrEqualPatchOnly,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Num: 1},
							Patch: NumConstraint{Num: 1},
						},
					},
				},
			},
			``,
		},
		{
			"~1.1",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpGreaterThanOrEqualPatchOnly,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Num: 1},
							Patch: NumConstraint{Num: 0},
						},
					},
				},
			},
			``,
		},
		{
			"~1",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpGreaterThanOrEqualMinorOnly,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Num: 0},
							Patch: NumConstraint{Num: 0},
						},
					},
				},
			},
			``,
		},
		{
			"^1.1.1",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpGreaterThanOrEqualMinorOnly,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Num: 1},
							Patch: NumConstraint{Num: 1},
						},
					},
				},
			},
			``,
		},
		{
			"^1.1",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpGreaterThanOrEqualMinorOnly,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Num: 1},
							Patch: NumConstraint{Num: 0},
						},
					},
				},
			},
			``,
		},
		{
			"^0.1",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpGreaterThanOrEqualPatchOnly,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 0},
							Minor: NumConstraint{Num: 1},
							Patch: NumConstraint{Num: 0},
						},
					},
				},
			},
			``,
		},
		{
			"^1",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpGreaterThanOrEqualMinorOnly,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Num: 0},
							Patch: NumConstraint{Num: 0},
						},
					},
				},
			},
			``,
		},
		{
			"=1.1.1",
			UnionSpec{
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
			},
			``,
		},
		{
			"!1.1.1",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpNotEqual,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Num: 1},
							Patch: NumConstraint{Num: 1},
						},
					},
				},
			},
			``,
		},
		{
			"1.*.*",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpMatch,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Unconstrained: true},
							Patch: NumConstraint{Unconstrained: true},
						},
					},
				},
			},
			``,
		},
		{
			"=1.*.*",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpEqual,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Num: 0},
							Patch: NumConstraint{Num: 0},
						},
					},
				},
			},
			``,
		},
		{
			"1.0.x",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpMatch,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Num: 0},
							Patch: NumConstraint{Unconstrained: true},
						},
					},
				},
			},
			``,
		},
		{
			"1.0.0 - 2.0.0",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpGreaterThanOrEqual,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Num: 0},
							Patch: NumConstraint{Num: 0},
						},
					},
					SelectionSpec{
						Operator: OpLessThanOrEqual,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 2},
							Minor: NumConstraint{Num: 0},
							Patch: NumConstraint{Num: 0},
						},
					},
				},
			},
			``,
		},
		{
			"1.*.* - 2.*.*",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpGreaterThanOrEqual,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Num: 0},
							Patch: NumConstraint{Num: 0},
						},
					},
					SelectionSpec{
						Operator: OpLessThan,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 3},
							Minor: NumConstraint{Num: 0},
							Patch: NumConstraint{Num: 0},
						},
					},
				},
			},
			``,
		},
		{
			">1.0.0 - 2.0.0",
			nil,
			`lower bound of range specified with "-" operator must be an exact version`,
		},
		{
			"1.0.0 - >2.0.0",
			nil,
			`upper bound of range specified with "-" operator must be an exact version`,
		},
		{
			">=1.0.0 <2.0.0",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpGreaterThanOrEqual,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Num: 0},
							Patch: NumConstraint{Num: 0},
						},
					},
					SelectionSpec{
						Operator: OpLessThan,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 2},
							Minor: NumConstraint{Num: 0},
							Patch: NumConstraint{Num: 0},
						},
					},
				},
			},
			``,
		},
		{
			">=1.0 <2 || 2.0-beta.1",
			UnionSpec{
				IntersectionSpec{
					SelectionSpec{
						Operator: OpGreaterThanOrEqual,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 1},
							Minor: NumConstraint{Num: 0},
							Patch: NumConstraint{Num: 0},
						},
					},
					SelectionSpec{
						Operator: OpLessThan,
						Boundary: VersionSpec{
							Major: NumConstraint{Num: 2},
							Minor: NumConstraint{Num: 0},
							Patch: NumConstraint{Num: 0},
						},
					},
				},
				IntersectionSpec{
					SelectionSpec{
						Operator: OpEqual,
						Boundary: VersionSpec{
							Major:      NumConstraint{Num: 2},
							Minor:      NumConstraint{Num: 0},
							Patch:      NumConstraint{Num: 0},
							Prerelease: "beta.1",
						},
					},
				},
			},
			``,
		},
		{
			"1.0.0, 2.0.0",
			nil,
			`commas are not needed to separate version selections; separate with spaces instead`,
		},
		{
			"= 1.1.1",
			nil,
			`no spaces allowed after operator "="`,
		},
		{
			"=  1.1.1",
			nil,
			`no spaces allowed after operator "="`,
		},
		{
			"garbage",
			nil,
			`the sequence "garbage" is not valid`,
		},
		{
			"&1.1.0",
			nil,
			`invalid constraint operator "&"`,
		},
		{
			"=v1.0.0",
			nil,
			`a "v" prefix should not be used when specifying versions`,
		},
		{
			"v1.0.0",
			nil,
			`a "v" prefix should not be used when specifying versions`,
		},
		{
			">=v1.0.0",
			nil,
			`a "v" prefix should not be used when specifying versions`,
		},
	}

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			got, err := Parse(test.Input)
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

			t.Logf("got: %s", spew.Sdump(got))
			t.Logf("want: %s", spew.Sdump(test.Want))

			for _, problem := range deep.Equal(got, test.Want) {
				t.Error(problem)
			}
		})
	}
}
