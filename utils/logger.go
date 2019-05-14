package utils

import (
	"time"

	"github.com/rivo/tview"
)

// FormatLog formats the message to display
func FormatLog(message string, mode string) string {
	messageEscaped := tview.Escape(message)
	dateTime := time.Now().Format(time.RFC3339)

	var value string
	switch mode {
	case "error":
		value = "[red]" + messageEscaped
	case "debug":
		value = "[orange]" + messageEscaped
	case "warn":
		value = "[yellow]" + messageEscaped
	case "data":
		return "[yellow]" + dateTime + " [" + BlueColorName + "]" + messageEscaped + "\r\n"
	default:
		value = "[white::-]" + messageEscaped
	}

	return "[yellow]" + dateTime + " " + value
}
