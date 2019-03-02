package utils

import (
	"io/ioutil"
	"os"

	"github.com/atotto/clipboard"
)

// WriteToClipboard writes data to the clipboard (Ctrl+C for example)
func WriteToClipboard(data string) {
	clipboard.WriteAll(data)
}

// ReadFromClipboard reads data from the clipboard
func ReadFromClipboard() (string, error) {
	return clipboard.ReadAll()
}

// GetByteFromPathFileName dffd
func GetByteFromPathFileName(pathFileName string, logger func(message string, mode string)) []byte {
	file, err := os.Open(pathFileName)
	if err != nil {
		logger("Not possible to load "+pathFileName+" file.", "error")
		return []byte(`{}`)
	}

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		logger("Not possible to read "+pathFileName+" file.", "error")
		return []byte(`{}`)
	}

	defer file.Close()

	return byteValue
}
