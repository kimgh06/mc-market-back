package utilities

import "cmp"

func Clamp[T cmp.Ordered](x T, minimum T, maximum T) T {
	return min(max(x, minimum), maximum)
}
