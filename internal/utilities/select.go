package utilities

// Select value depending on given boolean
func SelectB[T any](b bool, ifTrue T, ifFalse T) T {
	if b {
		return ifTrue
	} else {
		return ifFalse
	}
}
