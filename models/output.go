package models

import (
	"errors"
	"sort"
)

// Output struct corresponds to serialize and deserialize json app file
type Output struct {
	Data    []MakeRequestData
	Config  Config
	Context Context
}

// AddOrReplace adds or replaces a MakeRequestData struct
func (out *Output) AddOrReplace(data MakeRequestData) {
	// Initialize with the updated data
	newData := []MakeRequestData{data}
	for _, value := range out.Data {
		// URL & Method are the primary key
		if value.URL != data.URL || value.Method != data.Method {
			newData = append(newData, value)
		}
	}
	out.Data = newData
}

// Find finds a MakeRequestData from "method"/"url"
func (out Output) Find(method string, url string) (MakeRequestData, error) {
	var find MakeRequestData
	for _, value := range out.Data {
		if value.URL.String() == url && value.Method.String() == method {
			find = value
			return find, nil
		}
	}
	return find, errors.New("'" + method + " " + url + "' value does not exist")
}

// SortDataByProject returns map of project / data sorted by project
func (out Output) SortDataByProject() map[string][]MakeRequestData {

	sorted := func(slice map[string][]MakeRequestData) map[string][]MakeRequestData {
		// Sort by keys
		var keys []string
		for k := range slice {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// Re-build map with sorted keys
		new := make(map[string][]MakeRequestData)
		for _, k := range keys {
			new[k] = slice[k]
		}
		return new
	}

	add := func(newValue MakeRequestData, slice []MakeRequestData) []MakeRequestData {
		if slice == nil {
			return []MakeRequestData{newValue}
		}
		return append([]MakeRequestData{newValue}, slice...)
	}

	new := make(map[string][]MakeRequestData)
	for _, data := range out.Data {
		if data.ProjectName == "" {
			new["."] = add(data, new["."])
		} else {
			new[data.ProjectName] = add(data, new[data.ProjectName])
		}
	}

	return sorted(new)
}
