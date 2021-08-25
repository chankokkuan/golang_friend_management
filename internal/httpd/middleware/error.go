package middleware

import (
	"fmt"
	"friend-management/internal/core/domain"
	"friend-management/internal/httpd/middleware/jwtauth"
	"friend-management/internal/log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var codeToStatus = map[string]int{
	domain.NotFound:         http.StatusNotFound,
	domain.Unknown:          http.StatusInternalServerError,
	domain.AccessDenied:     http.StatusForbidden,
	domain.VersionConflict:  http.StatusConflict,
	domain.IsFriend:         http.StatusBadRequest,
	jwtauth.NoToken:         http.StatusUnauthorized,
	jwtauth.BadToken:        http.StatusUnauthorized,
	jwtauth.ClaimsNotFound:  http.StatusUnauthorized,
	jwtauth.ClaimsIncorrect: http.StatusUnauthorized,
}

func RespondWithError(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		errorToReturn := c.Errors.Last()
		if errorToReturn != nil {
			span, _ := tracer.SpanFromContext(c.Request.Context())

			if err, ok := errorToReturn.Err.(*domain.CustomErr); ok {
				status, isListed := codeToStatus[err.Code()]
				if !isListed {
					logger.WithCtx(c.Request.Context()).Warn("msg", fmt.Sprintf("Unmapped error code: %v", err.Code()))
					status = http.StatusInternalServerError
				}

				if len(err.Details()) > 0 {
					c.JSON(status, gin.H{
						"code":     err.Code(),
						"message":  err.Error(),
						"details":  err.Details(),
						"trace_id": span.Context().TraceID(),
					})
				} else {
					c.JSON(status, gin.H{
						"code":     err.Code(),
						"message":  err.Error(),
						"trace_id": span.Context().TraceID(),
					})
				}
			} else { // If forget to disguise as custom error
				logger.WithCtx(c.Request.Context()).Warn("msg", "Undisguised error returned as 500", "error", errorToReturn)
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":     domain.Unknown,
					"message":  domain.ErrUnknown.Error(),
					"trace_id": span.Context().TraceID(),
				})
			}
		}
	}
}
