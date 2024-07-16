package commands

import "time"

// Request: A request object
type Request struct {
	Name            string
	Url             string
	Method          string
	GetParams       string
	Data            string
	Headers         string
	ResponseHeaders string
	RawResponseBody []byte
	ContentType     string
	Duration        time.Duration
	// Formatter       formatter.ResponseFormatter
}

func (r *Request) Send() {

}
