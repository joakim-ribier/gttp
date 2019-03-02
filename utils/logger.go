package utils

import "time"

// FormatLog formats the message to display
func FormatLog(message string, mode string) string {
	t := time.Now()
	dateTime := t.Format(time.RFC3339)
	var value string
	switch mode {
	case "error":
		value = "[red]" + message
	case "debug":
		value = "[orange]" + message
	case "warn":
		value = "[yellow]" + message
	case "data":
		return "[yellow]" + dateTime + " [blue]" + message + "\r\n"
	default:
		value = "[white::-]" + message
	}
	return "[yellow]" + dateTime + " " + value
}
