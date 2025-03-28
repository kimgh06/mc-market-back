package utilities

func RemoveDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	var list []T
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func Filter[T any](ts []T, f func(T) bool) []T {
	var out []T
	for i := range ts {
		if f(ts[i]) {
			out = append(out, ts[i])
		}
	}
	return out
}

func Walk[T any](ts []T, f func(T) []T) {
	for i := range ts {
		Walk[T](f(ts[i]), f)
	}
}

func Sum[T any](ts []T, f func(T) int) int {
	var sum = 0
	for i := range ts {
		sum += f(ts[i])
	}
	return sum
}

func Count[T any](ts []T, f func(T) bool) int {
	var satisfies = 0
	for i := range ts {
		if f(ts[i]) {
			satisfies++
		}
	}

	return satisfies
}

func CountNil[T any](ts []*T) int {
	return Count(ts, func(t *T) bool {
		return t == nil
	})
}

func CountNotNil[T any](ts []*T) int {
	return Count(ts, func(t *T) bool {
		return t != nil
	})
}

func FlattenArray[V any](arr [][]V) []V {
	var newArr []V
	for _, a := range arr {
		newArr = append(newArr, a...)
	}

	return newArr
}
