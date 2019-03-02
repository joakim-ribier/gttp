package httpclient

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/joakim-ribier/gttp/models/types"
)

// Call http method
func Call(method types.Method, url types.URL, contentType string, data []byte, headers map[string]string, logger func(message string, mode string)) (*HTTPClient, error) {
	return getJSON(method, url, contentType, data, headers, logger)
}

func getJSON(method types.Method, url types.URL, contentType string, data []byte, headers map[string]string, logger func(message string, mode string)) (*HTTPClient, error) {
	logger(method.String()+" "+url.String(), "debug")

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest(method.String(), url.String(), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", contentType)

	// Set HTTP header values
	for key, value := range headers {
		if !(strings.HasPrefix(key, "{") && strings.HasSuffix(key, "}")) {
			req.Header.Set(key, value)
		}
	}

	if err != nil {
		logger("Impossible to build the query.", "error")
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		logger("Impossible to execute the query.", "error")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger("Impossible to read the response body.", "error")
		return nil, err
	}

	httpClient := NewHTTPClient(resp, body).withHeaderData(logger)
	return httpClient, nil
}
