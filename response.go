package treblle

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

type ResponseInfo struct {
	Headers  json.RawMessage `json:"headers"`
	Code     int             `json:"code"`
	Size     int             `json:"size"`
	LoadTime float64         `json:"load_time"`
	Body     json.RawMessage `json:"body"`
	Errors   []ErrorInfo     `json:"errors"`
}

type ErrorInfo struct {
	Source  string `json:"source"`
	Type    string `json:"type"`
	Message string `json:"message"`
	File    string `json:"file"`
	Line    int    `json:"line"`
}

type BodyLogWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

// Extract information from Gin body and context
func getGinResponseInfo(blw *BodyLogWriter, c *gin.Context, startTime time.Time) ResponseInfo {
	defer dontPanic()
	responseBytes := blw.Body.Bytes()

	errInfo := ErrorInfo{}
	var body map[string]interface{}
	err := json.Unmarshal(responseBytes, &body)
	if err != nil {
		errInfo.Message = err.Error()
	}

	headers := make(map[string]string)
	for k := range blw.Header() {
		headers[k] = blw.Header().Get(k)
	}

	// headers to Json
	headersJson, err := json.Marshal(headers)
	if err != nil {
		errInfo.Message = err.Error()
	}

	// body to Json
	bodyJson, err := json.Marshal(body)
	if err != nil {
		errInfo.Message = err.Error()
	}

	r := ResponseInfo{
		Headers:  headersJson,
		Code:     c.Writer.Status(),
		Size:     len(responseBytes),
		LoadTime: float64(time.Since(startTime).Microseconds()),
		Body:     bodyJson,
		Errors:   make([]ErrorInfo, 0),
	}

	if err != nil {
		r.Errors = append(r.Errors, errInfo)
	}

	return r
}
