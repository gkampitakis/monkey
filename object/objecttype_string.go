// Code generated by "stringer -type=ObjectType"; DO NOT EDIT.

package object

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[INTEGER-0]
	_ = x[BOOLEAN-1]
	_ = x[NULL-2]
	_ = x[RETURN_VALUE-3]
	_ = x[ERROR_VALUE-4]
}

const _ObjectType_name = "INTEGERBOOLEANNULLRETURN_VALUEERROR_VALUE"

var _ObjectType_index = [...]uint8{0, 7, 14, 18, 30, 41}

func (i ObjectType) String() string {
	if i >= ObjectType(len(_ObjectType_index)-1) {
		return "ObjectType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ObjectType_name[_ObjectType_index[i]:_ObjectType_index[i+1]]
}
