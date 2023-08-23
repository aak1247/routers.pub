package utils

func Map[S any, T any](s []S, f func(int, S) T) []T {
	res := make([]T, len(s))
	for i, item := range s {
		res[i] = f(i, item)
	}
	return res
}
