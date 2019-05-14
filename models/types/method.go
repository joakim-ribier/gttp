package types

import (
	"strings"

	"github.com/joakim-ribier/gttp/utils"
)

// Method string type value
type Method string

// String returns string value
func (m Method) String() string {
	return string(m)
}

// Label returns string to display
func (m Method) Label() string {
	str := m.String()
	return str + strings.Repeat(" ", len("DELETE")-len(str))
}

// TreeColor returns foreground & background
func (m Method) TreeColor() string {
	switch m.String() {
	case "GET":
		return "[" + utils.BlueColorName + ":]"
	case "POST":
		return "[" + utils.GreenColorName + ":]"
	case "PUT":
		return "[orange:]"
	case "DELETE":
		return "[red:]"
	default:
		return ""
	}
}
