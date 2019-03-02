package types

import (
	"fmt"
)

// Method string type value
type Method string

// String returns string value
func (m Method) String() string {
	return fmt.Sprintf("%s", string(m))
}

// TreeColor returns foreground & background
func (m Method) TreeColor() string {
	switch m.String() {
	case "GET":
		return "[white:blue]"
	case "POST":
		return "[white:green]"
	case "PUT":
		return "[white:orange]"
	default:
		return ""
	}
}
