package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

// Request: A request object
type Request struct {
	Name         string                 `json:"name"`
	Url          string                 `json:"url"`
	Method       string                 `json:"method"`
	ContentType  string                 `json:"contentType,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
	ResponseBody string                 `json:"rawResponseBody,omitempty"`
	// ResponseHeaders string `json:"responseHeaders"`
	// Headers         string `json:"headers"`
	// RawResponseBody byte   `json:"rawResponseBody"`
	// Duration        string `json:"duration"`
	// Duration        time.Duration
	// Formatter       formatter.ResponseFormatter

	Log         *logrus.Entry
	HttpCommand *HttpCommand
}

func (r *Request) Send() error {
    var request *http.Request
	var err error

	if r.Data != nil {
		jsonData, err := json.Marshal(r.Data)
		if err != nil {
			r.Log.Error("Error marshaling data: ", err)
		}
		request, err = http.NewRequest(r.Method, r.Url, bytes.NewBuffer(jsonData))
		request.Header.Set("Content-Type", r.ContentType)
	} else {
		request, err = http.NewRequest(r.Method, r.Url, nil)
	}
	if err != nil {
		return err
	}
	// request.Header.Set("do headers here")
	r.Log.Info("Sending request to: ", request.URL)
	response, err := r.HttpCommand.Client.Do(request)
	if err != nil {
		r.Log.Error("Request failed: ", request.URL, err)
		return err
	}
	responseBody, error := io.ReadAll(response.Body)

	if error != nil {
		fmt.Println(error)
	}

	formattedData := formatJSON(responseBody)
	r.ResponseBody = formattedData
	r.Log.Info("Response received: ", response.Status, r.ResponseBody)

	return nil
}

// TODO: Move to Utils
// function to format JSON data
func formatJSON(data []byte) string {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", " ")

	if err != nil {
		fmt.Println(err)
	}

	d := out.Bytes()
	return string(d)
}

func (r *Request) GetParams() (url.Values, error) {
	u, err := url.Parse(r.Url)
	if err != nil {
		return nil, err
	}
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}

	return m, nil
}
