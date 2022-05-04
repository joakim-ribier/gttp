package core

import "sort"

// StringMap map[string]string type
type StringMap map[string]string

// ToSliceOfKeys converts map[string]string to []string withs keys
func (sMap StringMap) ToSortedKeys() []string {
	tab := []string{}
	for key := range sMap {
		tab = append(tab, key)
	}
	sort.Strings(tab)
	return tab
}

// ReplaceContext replaces all intial values by the context values
func (sMap StringMap) ReplaceContext(mapKeysValues map[string]string) StringMap {
	new := make(map[string]string)
	for key, value := range sMap {
		new[key] = value
		if val, ok := mapKeysValues[value]; ok {
			new[key] = val
		}
	}
	return new
}
