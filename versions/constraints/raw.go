package constraints

//go:generate ragel -G1 -Z raw_scan.rl
//go:generate gofmt -w raw_scan.go

// rawConstraint is a tokenization of a constraint string, used internally
// as the first layer of parsing.
type rawConstraint struct {
	op    string
	sep   string
	nums  [3]string
	numCt int
	pre   string
	meta  string
}
