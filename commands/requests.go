package commands

import (
	"net/url"
)

// Request: A request object
type Request struct {
	Name      string `json:"name"`
	Url       string `json:"url"`
	Method    string `json:"method"`
	// Data            string `json:"data"`
	// Headers         string `json:"headers"`
	// ResponseHeaders string `json:"responseHeaders"`
	// RawResponseBody byte `json:"rawResponseBody"`
	// ContentType     string `json:"contentType"`
	// Duration        string `json:"duration"`
	// Duration        time.Duration
	// Formatter       formatter.ResponseFormatter
}

func (r *Request) Send() {

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
