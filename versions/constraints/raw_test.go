package constraints

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestScanConstraints(t *testing.T) {
	tests := []struct {
		Input      string
		Want       rawConstraint
		WantRemain string
	}{
		{
			"",
			rawConstraint{},
			"",
		},
		{
			"garbage",
			rawConstraint{},
			"garbage",
		},
		{
			"1.0.0",
			rawConstraint{
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
			},
			"",
		},
		{
			"1.0",
			rawConstraint{
				nums:  [...]string{"1", "0", ""},
				numCt: 2,
			},
			"",
		},
		{
			"1",
			rawConstraint{
				nums:  [...]string{"1", "", ""},
				numCt: 1,
			},
			"",
		},
		{
			"10.0.0",
			rawConstraint{
				nums:  [...]string{"10", "0", "0"},
				numCt: 3,
			},
			"",
		},
		{
			"10.0.0.0",
			rawConstraint{
				nums:  [...]string{"10", "0", "0"},
				numCt: 4,
			},
			"",
		},
		{
			"*",
			rawConstraint{
				nums:  [...]string{"*", "", ""},
				numCt: 1,
			},
			"",
		},
		{
			"*.*",
			rawConstraint{
				nums:  [...]string{"*", "*", ""},
				numCt: 2,
			},
			"",
		},
		{
			"*.*.*",
			rawConstraint{
				nums:  [...]string{"*", "*", "*"},
				numCt: 3,
			},
			"",
		},
		{
			"1.0.*",
			rawConstraint{
				nums:  [...]string{"1", "0", "*"},
				numCt: 3,
			},
			"",
		},
		{
			"x",
			rawConstraint{
				nums:  [...]string{"x", "", ""},
				numCt: 1,
			},
			"",
		},
		{
			"x.x",
			rawConstraint{
				nums:  [...]string{"x", "x", ""},
				numCt: 2,
			},
			"",
		},
		{
			"x.x.x",
			rawConstraint{
				nums:  [...]string{"x", "x", "x"},
				numCt: 3,
			},
			"",
		},
		{
			"1.0.x",
			rawConstraint{
				nums:  [...]string{"1", "0", "x"},
				numCt: 3,
			},
			"",
		},
		{
			"1.0.0-beta1",
			rawConstraint{
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
				pre:   "beta1",
			},
			"",
		},
		{
			"1.0.0+abc123",
			rawConstraint{
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
				meta:  "abc123",
			},
			"",
		},
		{
			"1.0.0-beta1+abc123",
			rawConstraint{
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
				pre:   "beta1",
				meta:  "abc123",
			},
			"",
		},
		{
			"1.0.0garbage",
			rawConstraint{
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
			},
			"garbage",
		},

		// Rubygems-style operators
		{
			">= 1.0.0",
			rawConstraint{
				op:    ">=",
				sep:   " ",
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
			},
			"",
		},
		{
			"<= 1.0.0",
			rawConstraint{
				op:    "<=",
				sep:   " ",
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
			},
			"",
		},
		{
			"> 1.0.0",
			rawConstraint{
				op:    ">",
				sep:   " ",
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
			},
			"",
		},
		{
			"< 1.0.0",
			rawConstraint{
				op:    "<",
				sep:   " ",
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
			},
			"",
		},
		{
			"= 1.0.0",
			rawConstraint{
				op:    "=",
				sep:   " ",
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
			},
			"",
		},
		{
			"!= 1.0.0",
			rawConstraint{
				op:    "!=",
				sep:   " ",
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
			},
			"",
		},
		{
			"~> 1.0.0",
			rawConstraint{
				op:    "~>",
				sep:   " ",
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
			},
			"",
		},
		{
			// comma separated, as sometimes seen in ruby-ish tools
			"1.0.0, 2.0.0",
			rawConstraint{
				op:    "",
				sep:   "",
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
			},
			", 2.0.0",
		},

		// npm-style operators
		{
			"<1.0.0",
			rawConstraint{
				op:    "<",
				sep:   "",
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
			},
			"",
		},
		{
			"<=1.0.0",
			rawConstraint{
				op:    "<=",
				sep:   "",
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
			},
			"",
		},
		{
			">1.0.0",
			rawConstraint{
				op:    ">",
				sep:   "",
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
			},
			"",
		},
		{
			">=1.0.0",
			rawConstraint{
				op:    ">=",
				sep:   "",
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
			},
			"",
		},
		{
			"=1.0.0",
			rawConstraint{
				op:    "=",
				sep:   "",
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
			},
			"",
		},
		{
			"~1.0.0",
			rawConstraint{
				op:    "~",
				sep:   "",
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
			},
			"",
		},
		{
			"^1.0.0",
			rawConstraint{
				op:    "^",
				sep:   "",
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
			},
			"",
		},
		{
			// npm-style range operator
			"1.0.0 - 2.0.0",
			rawConstraint{
				op:    "",
				sep:   "",
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
			},
			" - 2.0.0",
		},
		{
			// npm-style "or" operator
			"1.0.0 || 2.0.0",
			rawConstraint{
				op:    "",
				sep:   "",
				nums:  [...]string{"1", "0", "0"},
				numCt: 3,
			},
			" || 2.0.0",
		},
	}

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			got, remain := scanConstraint(test.Input)
			if remain != test.WantRemain {
				t.Errorf("wrong remain\ngot:  %q\nwant: %q", remain, test.WantRemain)
			}

			if diff := pretty.Compare(test.Want, got); diff != "" {
				t.Errorf("wrong result\n%s", diff)
			}
		})
	}
}
