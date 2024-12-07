//line raw_scan.rl:1
// This file is generated from raw_scan.rl. DO NOT EDIT.

//line raw_scan.rl:5

package constraints

//line raw_scan.go:9
var _scan_actions []byte = []byte{
	0, 1, 1, 1, 4, 1, 7, 2, 0,
	1, 2, 2, 1, 2, 3, 1, 2,
	3, 7, 2, 4, 1, 2, 4, 7,
	2, 5, 1, 2, 5, 7, 2, 6,
	7, 3, 0, 1, 2, 3, 2, 1,
	3, 3, 2, 3, 7, 4, 0, 1,
	2, 3,
}

var _scan_key_offsets []byte = []byte{
	0, 0, 14, 28, 37, 45, 53, 58,
	61, 69, 78,
}

var _scan_trans_keys []byte = []byte{
	32, 42, 46, 88, 118, 120, 9, 13,
	48, 57, 65, 90, 97, 122, 32, 42,
	46, 88, 118, 120, 9, 13, 48, 57,
	65, 90, 97, 122, 32, 42, 88, 118,
	120, 9, 13, 48, 57, 45, 46, 48,
	57, 65, 90, 97, 122, 45, 46, 48,
	57, 65, 90, 97, 122, 42, 88, 120,
	48, 57, 43, 45, 46, 45, 46, 48,
	57, 65, 90, 97, 122, 43, 45, 46,
	48, 57, 65, 90, 97, 122, 43, 45,
	46, 48, 57,
}

var _scan_single_lengths []byte = []byte{
	0, 6, 6, 5, 0, 0, 3, 3,
	0, 1, 3,
}

var _scan_range_lengths []byte = []byte{
	0, 4, 4, 2, 4, 4, 1, 0,
	4, 4, 1,
}

var _scan_index_offsets []byte = []byte{
	0, 0, 11, 22, 30, 35, 40, 45,
	49, 54, 60,
}

var _scan_indicies []byte = []byte{
	1, 2, 3, 2, 1, 2, 1, 4,
	3, 3, 0, 6, 7, 3, 7, 6,
	7, 6, 8, 3, 3, 5, 10, 11,
	11, 10, 11, 10, 12, 9, 14, 14,
	14, 14, 13, 15, 15, 15, 15, 13,
	16, 16, 16, 17, 13, 19, 20, 21,
	18, 14, 14, 14, 14, 22, 24, 15,
	15, 15, 15, 23, 19, 20, 21, 25,
	18,
}

var _scan_trans_targs []byte = []byte{
	2, 3, 7, 0, 10, 2, 3, 7,
	10, 0, 3, 7, 10, 0, 8, 9,
	7, 10, 0, 4, 5, 6, 0, 0,
	4, 10,
}

var _scan_trans_actions []byte = []byte{
	7, 34, 46, 42, 46, 0, 10, 38,
	38, 16, 0, 13, 13, 5, 0, 0,
	1, 1, 22, 19, 19, 3, 31, 28,
	25, 0,
}

var _scan_eof_actions []byte = []byte{
	0, 42, 42, 16, 5, 5, 5, 22,
	31, 28, 22,
}

const scan_start int = 1
const scan_first_final int = 7
const scan_error int = 0

const scan_en_main int = 1

//line raw_scan.rl:11

func scanConstraint(data string) (rawConstraint, string) {
	var constraint rawConstraint
	var numIdx int
	var extra string

	// Ragel state
	p := 0          // "Pointer" into data
	pe := len(data) // End-of-data "pointer"
	cs := 0         // constraint state (will be initialized by ragel-generated code)
	ts := 0
	te := 0
	eof := pe

	// Keep Go compiler happy even if generated code doesn't use these
	_ = ts
	_ = te
	_ = eof

//line raw_scan.go:112
	{
		cs = scan_start
	}

//line raw_scan.go:116
	{
		var _klen int
		var _trans int
		var _acts int
		var _nacts uint
		var _keys int
		if p == pe {
			goto _test_eof
		}
		if cs == 0 {
			goto _out
		}
	_resume:
		_keys = int(_scan_key_offsets[cs])
		_trans = int(_scan_index_offsets[cs])

		_klen = int(_scan_single_lengths[cs])
		if _klen > 0 {
			_lower := int(_keys)
			var _mid int
			_upper := int(_keys + _klen - 1)
			for {
				if _upper < _lower {
					break
				}

				_mid = _lower + ((_upper - _lower) >> 1)
				switch {
				case data[p] < _scan_trans_keys[_mid]:
					_upper = _mid - 1
				case data[p] > _scan_trans_keys[_mid]:
					_lower = _mid + 1
				default:
					_trans += int(_mid - int(_keys))
					goto _match
				}
			}
			_keys += _klen
			_trans += _klen
		}

		_klen = int(_scan_range_lengths[cs])
		if _klen > 0 {
			_lower := int(_keys)
			var _mid int
			_upper := int(_keys + (_klen << 1) - 2)
			for {
				if _upper < _lower {
					break
				}

				_mid = _lower + (((_upper - _lower) >> 1) & ^1)
				switch {
				case data[p] < _scan_trans_keys[_mid]:
					_upper = _mid - 2
				case data[p] > _scan_trans_keys[_mid+1]:
					_lower = _mid + 2
				default:
					_trans += int((_mid - int(_keys)) >> 1)
					goto _match
				}
			}
			_trans += _klen
		}

	_match:
		_trans = int(_scan_indicies[_trans])
		cs = int(_scan_trans_targs[_trans])

		if _scan_trans_actions[_trans] == 0 {
			goto _again
		}

		_acts = int(_scan_trans_actions[_trans])
		_nacts = uint(_scan_actions[_acts])
		_acts++
		for ; _nacts > 0; _nacts-- {
			_acts++
			switch _scan_actions[_acts-1] {
			case 0:
//line raw_scan.rl:33

				numIdx = 0
				constraint = rawConstraint{}

			case 1:
//line raw_scan.rl:38

				ts = p

			case 2:
//line raw_scan.rl:42

				te = p
				constraint.op = data[ts:p]

			case 3:
//line raw_scan.rl:47

				te = p
				constraint.sep = data[ts:p]

			case 4:
//line raw_scan.rl:52

				te = p
				constraint.numCt++
				if numIdx < len(constraint.nums) {
					constraint.nums[numIdx] = data[ts:p]
					numIdx++
				}

			case 5:
//line raw_scan.rl:61

				te = p
				constraint.pre = data[ts+1 : p]

			case 6:
//line raw_scan.rl:66

				te = p
				constraint.meta = data[ts+1 : p]

			case 7:
//line raw_scan.rl:71

				extra = data[p:]

//line raw_scan.go:245
			}
		}

	_again:
		if cs == 0 {
			goto _out
		}
		p++
		if p != pe {
			goto _resume
		}
	_test_eof:
		{
		}
		if p == eof {
			__acts := _scan_eof_actions[cs]
			__nacts := uint(_scan_actions[__acts])
			__acts++
			for ; __nacts > 0; __nacts-- {
				__acts++
				switch _scan_actions[__acts-1] {
				case 2:
//line raw_scan.rl:42

					te = p
					constraint.op = data[ts:p]

				case 3:
//line raw_scan.rl:47

					te = p
					constraint.sep = data[ts:p]

				case 4:
//line raw_scan.rl:52

					te = p
					constraint.numCt++
					if numIdx < len(constraint.nums) {
						constraint.nums[numIdx] = data[ts:p]
						numIdx++
					}

				case 5:
//line raw_scan.rl:61

					te = p
					constraint.pre = data[ts+1 : p]

				case 6:
//line raw_scan.rl:66

					te = p
					constraint.meta = data[ts+1 : p]

				case 7:
//line raw_scan.rl:71

					extra = data[p:]

//line raw_scan.go:303
				}
			}
		}

	_out:
		{
		}
	}

//line raw_scan.rl:92

	return constraint, extra
}
