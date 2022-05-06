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

// EmptyMakeRequestData creates an empty new MakeRequestData struct
func EmptyMakeRequestData() MakeRequestData {
	return NewMakeRequestData("GET", "", make(core.StringMap), "", "application/json", "", "")
}

// SimpleMakeRequestData creates a simple new MakeRequestData
func SimpleMakeRequestData(
	method string,
	url string,
	projectName string,
	alias string) MakeRequestData {

	return NewMakeRequestData(method, url, make(core.StringMap), "", "application/json", projectName, alias)
}

// NewMakeRequestData creates a new MakeRequestData
func NewMakeRequestData(
	method string,
	url string,
	header core.StringMap,
	body string,
	contentType string,
	projectName string,
	alias string) MakeRequestData {

	return MakeRequestData{
		Method:                   types.Method(method),
		URL:                      types.URL(url),
		MapRequestHeaderKeyValue: header,
		Body:                     body,
		ContentType:              contentType,
		ProjectName:              projectName,
		Alias:                    alias,
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

// ToLog builds request to str message to be logged
func (m MakeRequestData) ToLog(url types.URL) string {
	var sb strings.Builder
	sb.WriteString(string(m.Method))
	sb.WriteString(" ")
	sb.WriteString(string(url))
	return sb.String()
}
