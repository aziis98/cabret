package util

func CloneMap[K comparable, V any](m1 map[K]V) map[K]V {
	m2 := map[K]V{}
	for k, v := range m1 {
		m2[k] = v
	}
	return m2
}
