package actions

import "github.com/joakim-ribier/gttp/httpclient"

type MakeRequestAction struct {
	DisplayResponse     func(client *httpclient.HTTPClient, data string)
	DisplayErrorRequest func(message string, mode string)
}

func NewMakeRequestAction(
	displayResponse func(client *httpclient.HTTPClient, data string),
	displayErrorRequest func(message string, mode string)) *MakeRequestAction {

	return &MakeRequestAction{
		DisplayResponse:     displayResponse,
		DisplayErrorRequest: displayErrorRequest,
	}
}
