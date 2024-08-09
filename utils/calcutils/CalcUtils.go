package calcutils

// Intersection 交集
func Intersection(set1, set2 map[string]struct{}) map[string]struct{} {
	result := make(map[string]struct{})
	for k := range set1 {
		if _, ok := set2[k]; ok {
			result[k] = struct{}{}
		}
	}
	return result
}

// Union 并集
func Union(set1, set2 map[string]struct{}) map[string]struct{} {
	result := make(map[string]struct{})
	for k := range set1 {
		result[k] = struct{}{}
	}
	for k := range set2 {
		result[k] = struct{}{}
	}
	return result
}

// Difference 差集
func Difference(set1, set2 map[string]struct{}) map[string]struct{} {
	result := make(map[string]struct{})
	for k := range set1 {
		if _, ok := set2[k]; !ok {
			result[k] = struct{}{}
		}
	}
	return result
}
