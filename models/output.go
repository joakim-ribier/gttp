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

// Remove removes MakeRequestData struct
func (out *Output) Remove(data MakeRequestData) {
	newData := []MakeRequestData{}
	for _, value := range out.Data {
		// URL & Method are the primary key
		if !(value.URL == data.URL && value.Method == data.Method) {
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

// SortDataAPIsByProjectName filters data APIs by project name and sorts by them
func (out Output) SortDataAPIsByProjectName() ([]string, map[string][]MakeRequestData) {

	getSortedKeys := func(values map[string][]MakeRequestData) []string {
		var keys []string
		for k := range values {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return keys
	}

	filterDataAPIsByProjectName := func(data []MakeRequestData) map[string][]MakeRequestData {

		add := func(newValue MakeRequestData, slice []MakeRequestData) []MakeRequestData {
			if slice == nil {
				return []MakeRequestData{newValue}
			}
			return append([]MakeRequestData{newValue}, slice...)
		}

		mapByProjectName := make(map[string][]MakeRequestData)
		for _, data := range data {
			if data.ProjectName == "" {
				mapByProjectName["."] = add(data, mapByProjectName["."])
			} else {
				mapByProjectName[data.ProjectName] = add(data, mapByProjectName[data.ProjectName])
			}
		}
		return mapByProjectName
	}

	dataAPIsByProjectName := filterDataAPIsByProjectName(out.Data)

	return getSortedKeys(dataAPIsByProjectName), dataAPIsByProjectName
}

func (out Output) UpdateMakeRequestData(values []MakeRequestData) Output {
	return Output{values, out.Config, out.Context}
}
