package arrayutil

func Reverse[T any](arr []T) []T {
	c := make([]T, len(arr))
	copy(c, arr)
	for i, j := 0, len(c)-1; i < j; i, j = i+1, j-1 {
		c[i], c[j] = c[j], c[i]
	}
	return c
}
