package constraints

import (
	"testing"

	"github.com/go-test/deep"
)

func TestParseExactVersion(t *testing.T) {
	tests := []struct {
		Input   string
		Want    VersionSpec
		WantErr string
	}{
		{
			"",
			VersionSpec{},
			"empty specification",
		},
		{
			"1",
			VersionSpec{
				Major: NumConstraint{Num: 1},
				Minor: NumConstraint{Num: 0},
				Patch: NumConstraint{Num: 0},
			},
			"",
		},
		{
			"1.1",
			VersionSpec{
				Major: NumConstraint{Num: 1},
				Minor: NumConstraint{Num: 1},
				Patch: NumConstraint{Num: 0},
			},
			"",
		},
		{
			"1.1.1",
			VersionSpec{
				Major: NumConstraint{Num: 1},
				Minor: NumConstraint{Num: 1},
				Patch: NumConstraint{Num: 1},
			},
			"",
		},
		{
			"1.0.0.0",
			VersionSpec{},
			"too many numbered portions; only three are allowed (major, minor, patch)",
		},
		{
			"v1.0.0",
			VersionSpec{},
			`a "v" prefix should not be used`,
		},
		{
			"1.0.0-beta2",
			VersionSpec{
				Major:      NumConstraint{Num: 1},
				Minor:      NumConstraint{Num: 0},
				Patch:      NumConstraint{Num: 0},
				Prerelease: "beta2",
			},
			"",
		},
		{
			"1.0-beta2",
			VersionSpec{
				Major:      NumConstraint{Num: 1},
				Minor:      NumConstraint{Num: 0},
				Patch:      NumConstraint{Num: 0},
				Prerelease: "beta2",
			},
			"",
		},
		{
			"1.0.0-beta.2",
			VersionSpec{
				Major:      NumConstraint{Num: 1},
				Minor:      NumConstraint{Num: 0},
				Patch:      NumConstraint{Num: 0},
				Prerelease: "beta.2",
			},
			"",
		},
		{
			"1.0.0+foo",
			VersionSpec{
				Major:    NumConstraint{Num: 1},
				Minor:    NumConstraint{Num: 0},
				Patch:    NumConstraint{Num: 0},
				Metadata: "foo",
			},
			"",
		},
		{
			"1.0.0+foo.bar",
			VersionSpec{
				Major:    NumConstraint{Num: 1},
				Minor:    NumConstraint{Num: 0},
				Patch:    NumConstraint{Num: 0},
				Metadata: "foo.bar",
			},
			"",
		},
		{
			"1.0.0-beta1+foo.bar",
			VersionSpec{
				Major:      NumConstraint{Num: 1},
				Minor:      NumConstraint{Num: 0},
				Patch:      NumConstraint{Num: 0},
				Prerelease: "beta1",
				Metadata:   "foo.bar",
			},
			"",
		},
		{
			"> 1.1.1",
			VersionSpec{},
			`can't use constraint operator ">"; an exact version is required`,
		},
		{
			"garbage",
			VersionSpec{},
			`invalid specification; required format is three positive integers separated by periods`,
		},
		{
			"& 1.1.0",
			VersionSpec{},
			`invalid sequence "&" at start of specification`,
		},
		{
			"1.*.*",
			VersionSpec{},
			`can't use wildcard for minor number; an exact version is required`,
		},
		{
			"1.0.x",
			VersionSpec{},
			`can't use wildcard for patch number; an exact version is required`,
		},
		{
			"1.0 || 2.0",
			VersionSpec{},
			`can't specify multiple versions; a single exact version is required`,
		},
		{
			"1.0.0, 2.0.0",
			VersionSpec{},
			`can't specify multiple versions; a single exact version is required`,
		},
		{
			"1.0.0 - 2.0.0",
			VersionSpec{},
			`can't specify version range; a single exact version is required`,
		},
	}

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			got, err := ParseExactVersion(test.Input)
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
