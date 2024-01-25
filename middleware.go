package treblle

import (
	"bytes"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	treblleVersion = "0.7.2"
	sdkName        = "go"
)

func GinMiddleware() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		startTime := time.Now()

		requestInfo, errReqInfo := getRequestInfo(gctx.Request, startTime)

		blw := &BodyLogWriter{Body: bytes.NewBufferString(""), ResponseWriter: gctx.Writer}
		gctx.Writer = blw

		gctx.Next()

		if !errors.Is(errReqInfo, ErrNotJson) {
			ti := MetaData{
				ApiKey:    Config.APIKey,
				ProjectID: Config.ProjectID,
				Version:   treblleVersion,
				Sdk:       sdkName,
				Data: DataInfo{
					Server:   Config.serverInfo,
					Language: Config.languageInfo,
					Request:  requestInfo,
					Response: getGinResponseInfo(blw, gctx.Copy(), startTime),
				},
			}

			// don't block execution while sending data to Treblle
			go sendToTreblle(ti)
		}
	}
}

// If anything happens to go wrong inside one of treblle-go internals, recover from panic and continue
func dontPanic() {
	if err := recover(); err != nil {
		log.Printf("treblle-go panic: %s", err)
	}
}
