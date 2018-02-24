package versions

import (
	"sort"
)

var _ sort.Interface = List(nil)
