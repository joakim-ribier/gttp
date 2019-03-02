package httpclient

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
)

// HTTPClient is object which contains *http data.
type HTTPClient struct {
	Request         *HTTPRequestClient
	Response        *HTTPResponseClient
	Body            []byte
	HeadersRequest  map[string]string
	HeadersResponse map[string]string
}

// HTTPRequestClient struct
type HTTPRequestClient struct {
	Host   string
	Method string
	HTTP   string
	URL    string
	Body   string
}

// HTTPResponseClient struct
type HTTPResponseClient struct {
	Response       *http.Response
	Status         string
	StatusCode     string
	HTTP           string
	Contentlength  string
	ContentType    string
	Date           string
	Referrerpolicy string
	Connection     string
}

// NewHTTPClient returns an instance of HTTPClient object which contains *httpResponse and body.
func NewHTTPClient(response *http.Response, data []byte) *HTTPClient {
	return &HTTPClient{
		Request:         newHTTPRequestClient(response),
		Response:        newHTTPResponseClient(response),
		Body:            data,
		HeadersRequest:  make(map[string]string),
		HeadersResponse: make(map[string]string),
	}
}

func newHTTPRequestClient(response *http.Response) *HTTPRequestClient {
	body := func(request *http.Request) string {
		body, _ := request.GetBody()
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(body)
		if err != nil {
			return "NOTE: binary data not shown in terminal"
		}
		return buf.String()
	}

	return &HTTPRequestClient{
		Host:   response.Request.URL.Host,
		Method: response.Request.Method,
		HTTP:   strconv.Itoa(response.Request.ProtoMajor) + "." + strconv.Itoa(response.Request.ProtoMinor),
		URL:    response.Request.URL.Path,
		Body:   body(response.Request),
	}
}

func newHTTPResponseClient(response *http.Response) *HTTPResponseClient {
	return &HTTPResponseClient{
		Response:       response,
		Status:         response.Status,
		StatusCode:     strconv.Itoa(response.StatusCode),
		HTTP:           strconv.Itoa(response.ProtoMajor) + "." + strconv.Itoa(response.ProtoMinor),
		Contentlength:  strconv.FormatInt(response.ContentLength, 10),
		ContentType:    response.Header.Get("content-type"),
		Date:           response.Header.Get("date"),
		Referrerpolicy: response.Header.Get("referrer-policy"),
		Connection:     response.Header.Get("connection"),
	}
}

func (client *HTTPClient) headerRequest(key string, value string) *HTTPClient {
	client.HeadersRequest[key] = value
	return client
}

func (client *HTTPClient) headerResponse(key string, value string) *HTTPClient {
	client.HeadersResponse[key] = value
	return client
}

func (client *HTTPClient) withHeaderData(logger func(message string, mode string)) *HTTPClient {
	response := client.Response.Response

	for k, v := range response.Request.Header {
		client = client.headerRequest(k, v[0])
	}

	for k, v := range response.Header {
		client = client.headerResponse(k, v[0])
	}

	logger("Request header: "+fmt.Sprintf("%s", response.Request.Header), "debug")
	logger("Response header: "+fmt.Sprintf("%s", response.Header), "debug")

	return client
}
