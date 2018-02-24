package versions

import (
	"encoding/json"
	"testing"

	"github.com/go-test/deep"
)

func TestVersionJSON(t *testing.T) {
	v := Version{
		Major:      1,
		Minor:      2,
		Patch:      3,
		Prerelease: "beta.1",
		Metadata:   "tci.12345",
	}

	j, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := string(j), `"1.2.3-beta.1+tci.12345"`; got != want {
		t.Errorf("wrong result\ngot:  %s\nwant: %s", got, want)
	}

	var v2 Version
	err = json.Unmarshal(j, &v2)
	if err != nil {
		t.Fatal(err)
	}

	for _, problem := range deep.Equal(v2, v) {
		t.Error(problem)
	}

	bad := []byte(`"garbage"`)
	err = json.Unmarshal(bad, &v2)
	var errText string
	if err != nil {
		errText = err.Error()
	}
	if got, want := errText, `invalid specification; required format is three positive integers separated by periods`; got != want {
		t.Errorf("wrong error for garbage\ngot:  %s\nwant: %s", got, want)
	}

}
