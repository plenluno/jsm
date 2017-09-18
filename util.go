package jsm

import (
	"fmt"
	"math"
)

const (
	intSize = 32 << (^uint(0) >> 63)
	maxInt  = 1<<(intSize-1) - 1
	minInt  = -1 << (intSize - 1)
)

func floatToInt(f float64) int {
	if math.IsNaN(f) {
		return 0
	}

	if math.IsInf(f, 0) {
		if math.Signbit(f) {
			return minInt
		}
		return maxInt
	}

	var sign float64
	if math.Signbit(f) {
		sign = -1.0
	} else {
		sign = 1.0
	}
	return int(sign * math.Floor(math.Abs(f)))
}

func floatToString(f float64) string {
	if math.IsInf(f, 0) {
		if math.Signbit(f) {
			return "-Infinity"
		}
		return "Infinity"
	}

	return fmt.Sprintf("%v", f)
}
