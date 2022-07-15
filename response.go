package treblle

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

type ResponseInfo struct {
	Headers  map[string]string      `json:"headers"`
	Code     int                    `json:"code"`
	Size     int                    `json:"size"`
	LoadTime float64                `json:"load_time"`
	Body     map[string]interface{} `json:"body"`
	Errors   []ErrorInfo            `json:"errors"`
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

func (w BodyLogWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Extract information from Gin body and context
func getGinResponseInfo(blw *BodyLogWriter, c *gin.Context, startTime time.Time) ResponseInfo {
	defer dontPanic()
	responseBytes := []byte(blw.Body.String())

	errInfo := ErrorInfo{}
	var body map[string]interface{}
	err := json.Unmarshal(responseBytes, &body)
	if err != nil {
		errInfo.Message = err.Error()
	}

	headers := make(map[string]string)
	for k, _ := range blw.Header() {
		headers[k] = blw.Header().Get(k)
	}

	r := ResponseInfo{
		Headers:  headers,
		Code:     c.Writer.Status(),
		Size:     len(responseBytes),
		LoadTime: float64(time.Since(startTime).Microseconds()),
		Body:     body,
		Errors:   make([]ErrorInfo, 0),
	}

	if err != nil {
		r.Errors = append(r.Errors, errInfo)
	}

	return r
}
