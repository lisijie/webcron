package models

func MakeIntSliceUnique(s []int) []int {
	rt := make([]int, 0)
	m := map[int]bool{}
	ok := false
	for _, v := range s {
		if _, ok = m[v]; !ok {
			rt = append(rt, v)
			m[v] = true
		}
	}
	return rt
}
