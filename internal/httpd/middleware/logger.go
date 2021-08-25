package middleware

import (
	"bytes"
	"friend-management/internal/log"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LogRequestAndResponse(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		reqBody, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

		// Attach logger
		span, _ := tracer.SpanFromContext(c.Request.Context())
		newCtx := logger.Attach(c.Request.Context(), "trace_id", span.Context().TraceID())
		c.Request = c.Request.WithContext(newCtx)

		c.Next()

		req := string(reqBody)
		if len(req) > 300 {
			req = req[:300] + "..."
		}

		var resp = blw.body.String()
		if len(resp) > 300 {
			resp = resp[:300] + "..."
		}

		logger.WithCtx(c.Request.Context()).Info(
			"method_path", c.Request.Method+" "+c.Request.URL.Path,
			"status_code", c.Writer.Status(),
			"request", req,
			"response", resp,
		)
	}
}
