package array

func Unique[T comparable](arr []T) []T {
	var unique []T
	m := make(map[T]bool)

	for _, v := range arr {
		if !m[v] {
			m[v] = true
			unique = append(unique, v)
		}
	}

	return unique

}
