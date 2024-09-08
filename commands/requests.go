package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/g-e-e-z/cucu/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Request: A request object
type Request struct {
	Uuid            string                 `json:"uuid"`
	Name            string                 `json:"name"`
	Url             string                 `json:"url"`
	Method          string                 `json:"method"`
	Headers         http.Header            `json:"headers,omitempty"`
	ContentType     string                 `json:"contentType,omitempty"`
	Data            map[string]interface{} `json:"data,omitempty"`
	Status          string
	ResponseBody    string
	ResponseHeaders http.Header
	Duration        time.Duration

	// Formatter       formatter.ResponseFormatter
	Hash string

	saved bool

	Log       *logrus.Entry
	Client    *http.Client
	OSCommand *OSCommand
	Modified  bool
}

func NewRequest(log *logrus.Entry, client *http.Client) *Request {
	request := &Request{
		Uuid:   uuid.New().String(),
		Name:   "NewReq",
		Url:    "placeholder url",
		Method: http.MethodGet,
		Log:    log,
		Client: client,
	}
	// This feels silly
	request.Hash = request.CreateHash()
	return request
}

func (r *Request) CheckModifed() {
	if r.Hash != r.CreateHash() || r.saved != true {
		r.Modified = true
	} else {
		r.Modified = false
	}
}

func (r *Request) CreateHash() string {
	return r.Name + r.Method + r.Url + r.DataToJSON() + r.HeadersToJSON()
}

func (r *Request) HeadersToJSON() string {
	bytes, err := json.Marshal(r.Headers)
	if err != nil {
		r.Log.WithError(err).Error("Failed to marshal JSON")
		return ""
	}
	return string(bytes)
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
		"headers":     r.Headers,
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
	if !r.Modified && r.saved {
		return nil
	}
	return r.OSCommand.SaveRequest(r)
}

func (r *Request) Delete(requests []*Request) error {
	return r.OSCommand.DeleteRequest(r, requests)
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
	if r.Headers != nil {
		request.Header = r.Headers
	} else {
		request.Header = http.Header{}
	}
	request.Header.Set("user-agent", "Cucu/0.0.1")
	r.Log.Info("Sending request to: ", request.URL)
	startTime := time.Now()
	response, err := r.Client.Do(request)
	if err != nil {
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

func (r *Request) GetRequestHeaders() (map[string]interface{}, error) {
	if r.Headers == nil {
		return nil, errors.New("request headers is empty")
	}

	result := make(map[string]interface{})
	for key, value := range r.Headers {
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
