package versions

// Make sure our set implementations all actually implement the interface
var _ setI = setBound{}
var _ setI = setExact{}
var _ setI = setExtreme(true)
var _ setI = setIntersection{}
var _ setI = setSubtract{}
var _ setI = setUnion{}
