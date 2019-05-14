package models

import (
	"strings"

	"github.com/joakim-ribier/gttp/core"
	"github.com/joakim-ribier/gttp/models/types"
	"github.com/joakim-ribier/gttp/utils"
)

// MakeRequestData reprensents a request structure
type MakeRequestData struct {
	Method                   types.Method
	URL                      types.URL
	MapRequestHeaderKeyValue core.StringMap
	Body                     string
	ContentType              string
	ProjectName              string
	Alias                    string
}

// NewMakeRequestData creates new MakeRequestData struct
func NewMakeRequestData() MakeRequestData {
	return MakeRequestData{
		Method:                   types.Method("GET"),
		URL:                      "",
		MapRequestHeaderKeyValue: make(core.StringMap),
		Body:                     "",
		ContentType:              "application/json",
		ProjectName:              "",
		Alias:                    "",
	}
}

func (m MakeRequestData) TreeFormat(pattern string) string {
	if pattern == "" {
		return m.Method.Label() + " " + m.URL.String()
	}
	value := strings.Replace(pattern, "{m}", m.Method.Label(), -1)
	if strings.Contains(value, "{a}|{u}") {
		if m.Alias != "" {
			value = strings.Replace(value, "{a}|{u}", m.Alias, -1)
		} else {
			value = strings.Replace(value, "{a}|{u}", m.URL.Base(), -1)
		}
	}
	if strings.Contains(value, "{a}|{url}") {
		if m.Alias != "" {
			value = strings.Replace(value, "{a}|{url}", m.Alias, -1)
		} else {
			value = strings.Replace(value, "{a}|{url}", m.URL.String(), -1)
		}
	}
	value = strings.Replace(value, "{color}", m.Method.TreeColor(), -1)
	value = strings.Replace(value, "{url}", m.URL.String(), -1)
	value = strings.Replace(value, "{u}", m.URL.Base(), -1)
	value = strings.Replace(value, "{backColor}", utils.BackColorName, -1)
	return value
}

// GetHTTPHeaderValues filters by HTTP request header params
func (m MakeRequestData) GetHTTPHeaderValues() core.StringMap {
	new := make(core.StringMap)
	for key, value := range m.MapRequestHeaderKeyValue {
		if !(strings.HasPrefix(key, "{") && strings.HasSuffix(key, "}")) {
			new[key] = value
		}
	}
	return new
}
