package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	// "go/token"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/g-e-e-z/cucu/utils"
	"github.com/sirupsen/logrus"
)

// Request: A request object
type Request struct {
	Uuid            string                 `json:"uuid"`
	Name            string                 `json:"name"`
	Url             string                 `json:"url"`
	Method          string                 `json:"method"`
	ContentType     string                 `json:"contentType,omitempty"`
	Data            map[string]interface{} `json:"data,omitempty"`
	Status          string
	ResponseBody    string
	ResponseHeaders http.Header
	Duration        time.Duration

	// Headers         string `json:"headers"`
	// Formatter       formatter.ResponseFormatter
	hash string

	Log         *logrus.Entry
	HttpCommand *HttpCommand
	Modified    bool
}

func (r *Request) CheckModifed() {
	if r.hash != r.createHash() {
		r.Modified = true
	} else {
		r.Modified = false
	}
}

func (r *Request) createHash() string {
	return r.Name + r.Method + r.Url + r.DataToJSON()
}

func (r *Request) DataToJSON() string {
    bytes, err := json.Marshal(r.Data)
	if err != nil {
		r.Log.WithError(err).Error("Failed to marshal JSON")
		return ""
	}
	return string(bytes)
}

func (r *Request) toJSON() string {
	// Create a map to hold the JSON representation
	jsonMap := map[string]interface{}{
		"uuid":        r.Uuid,
		"name":        r.Name,
		"url":         r.Url,
		"method":      r.Method,
		"contentType": r.ContentType,
		"data":        r.Data,
	}

	// Marshal the map into a JSON string
	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		r.Log.WithError(err).Error("Failed to marshal JSON")
		return ""
	}

	return string(jsonBytes)
}

func (r *Request) Save() error {
	if !r.Modified {
		return nil
	}
	return r.HttpCommand.SaveRequest(r)
}

func (r *Request) Delete() error {
	return nil //r.HttpCommand.DeleteRequest(r)
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
	startTime := time.Now()
	response, err := r.HttpCommand.Client.Do(request)
	if err != nil {
		// TODO: This handling is bad
		r.Log.Error("Request failed: ", request.URL, err)
		r.Status = "503 Service Unavailable"
		r.Duration = time.Since(startTime)
		r.ResponseBody = err.Error()
		return nil
	}
	r.Duration = time.Since(startTime)
	r.Status = response.Status
	r.ResponseHeaders = response.Header

	responseBody, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	contentType := r.ResponseHeaders.Get("Content-Type")
	var formattedData string
	if strings.Contains(contentType, "json") {
		formattedData = utils.FormatJSON(responseBody)
	}
	r.ResponseBody = formattedData
	r.Log.Info("Response received: ", response.Status, contentType, r.Duration)

	return nil
}

func (r *Request) GetData() (map[string]interface{}, error) {
	if r.Data == nil {
		return nil, errors.New("request data is empty")
	}

	result := make(map[string]interface{})
	for key, value := range r.Data {
		result[key] = value
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
