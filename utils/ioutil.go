package utils

import (
	"io/ioutil"

	"github.com/atotto/clipboard"
)

// WriteToClipboard writes data to the clipboard (Ctrl+C for example)
func WriteToClipboard(data string, logger func(message string, event string)) {
	if error := clipboard.WriteAll(data); error != nil {
		logger("Error to write data on clipboard.", "warn")
	}
}

// ReadFromClipboard reads data from the clipboard
func ReadFromClipboard() (string, error) {
	return clipboard.ReadAll()
}

// ReadFile reads the file @filename and return the content.
// If ioutil.ReadFile return an error, method logs and returns an empty cotent ([]byte).
func ReadFile(filename string, logger func(message string, mode string)) []byte {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logger("Reading data from '"+filename+"' file error...", "error")
		return []byte(`{}`)
	}
	return content
}
