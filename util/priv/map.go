package priv

import "strings"

// StrUniq returns reduced equivemnt str
//
func StrUniq(texts []string) []string {
	m := make(map[string]struct{})
	for _, t := range texts {
		m[t] = struct{}{}
	}

	uniq := []string{}
	for t := range m {
		uniq = append(uniq, t)
	}

	return uniq
}

// GetMapValue walks the dot-delimited `path` to return a nested map value, or nil.
//
func GetMapValue(m map[string]interface{}, path string) interface{} {
	var obj interface{} = m
	var val interface{}

	parts := strings.Split(path, ".")
	for _, p := range parts {
		if v, ok := obj.(map[string]interface{}); ok {
			obj = v[p]
			val = obj
		} else {
			return nil
		}
	}

	return val
}
