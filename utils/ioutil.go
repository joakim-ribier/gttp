package utils

import (
	"io/ioutil"
	"os"

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

// GetByteFromPathFileName dd
func GetByteFromPathFileName(pathFileName string, logger func(message string, mode string)) []byte {
	file, err := os.Open(pathFileName)
	if err != nil {
		logger("Error to load '"+pathFileName+"' file.", "error")
		return []byte(`{}`)
	}

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		logger("Error to read '"+pathFileName+"' file.", "error")
		return []byte(`{}`)
	}

	defer file.Close()

	return byteValue
}
