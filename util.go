package gopandas

import "time"

//StringArrayGetElse get string value from array or else
func StringArrayGetElse(array []string, idx int, elseValue string) string {
	if idx >= 0 && idx < len(array) {
		return array[idx]
	}
	return elseValue
}

//IntRange gives a int slice which start from start and end( not include )at end with step.
//If step doesn't have the same sign as the end-start, it will return a empty slice
func IntRange(start, end, step int) []int {
	if (end-start)*step <= 0 {
		return nil
	}
	ret := make([]int, 0, end-start/step)
	for i := start; i < end; i += step {
		ret = append(ret, i)
	}
	return ret
}

var (
	AncientTime = time.Unix(0, 0)
)
