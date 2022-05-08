package services

import (
	"encoding/json"
	"io/ioutil"

	"github.com/joakim-ribier/gttp/models"
	"github.com/joakim-ribier/gttp/utils"
)

type ApplicationDataService struct {
	Filename string
	Log      func(string, string)
}

// NewApplicationDataService constructs service which loads and saves the data from configuration file.
func NewApplicationDataService(filename string, log func(string, string)) *ApplicationDataService {
	return &ApplicationDataService{
		Filename: filename,
		Log:      log,
	}
}

// Load deserializes json file to a @models.Output.
func (s *ApplicationDataService) Load() models.Output {
	var value models.Output

	bytes := utils.ReadFile(s.Filename, s.Log)
	if error := json.Unmarshal(bytes, &value); error != nil {
		s.Log("Error to decode '"+s.Filename+"' json data file.", "error")
	}

	return value
}

// Save serializes @models.Output in the configuration app file.
func (s *ApplicationDataService) Save(value models.Output) {
	if json, error := json.Marshal(value); error != nil {
		s.Log("Encoding 'output' model error...", "error")
	} else {
		if error := ioutil.WriteFile(s.Filename, json, 0644); error != nil {
			s.Log("Writing data to '"+s.Filename+"' file error...", "error")
		}
	}
}
