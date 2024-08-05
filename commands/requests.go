package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/g-e-e-z/cucu/utils"
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
        // TODO: This handling is bad
		r.Log.Error("Request failed: ", request.URL, err)
        r.ResponseBody = err.Error()
	} else {
		responseBody, error := io.ReadAll(response.Body)
        defer response.Body.Close()

		if error != nil {
			fmt.Println(error)
		}

		formattedData := utils.FormatJSON(responseBody)
		r.ResponseBody = formattedData
		r.Log.Info("Response received: ", response.Status, r.ResponseBody)
	}

	return nil
}
func (r *Request) GetData() (map[string]string, error) {
	if r.Data == nil {
		return nil, errors.New("data is nil")
	}

    result := make(map[string]string)
	for key, value := range r.Data {
		result[key] = fmt.Sprintf("%v", value)
	}

	return result, nil
}
func (r *Request) GetParams() ([][]string, error) {
    params, err := utils.Parse(r.Url)
	if err != nil {
		return nil, err
	}

	return params, nil
}


