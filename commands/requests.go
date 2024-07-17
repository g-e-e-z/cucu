package commands

// import "time"

// Request: A request object
type Request struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	Method    string `json:"method"`
	// GetParams struct  `json:"getParams"`
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
