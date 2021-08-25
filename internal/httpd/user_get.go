package httpd

import (
	"friend-management/internal/httpd/middleware/jwtauth"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) UserGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Keys[jwtauth.JWTClaimsCtxKey].(jwtauth.JWTClaims)

		user, err := s.repo.GetUser(c.Request.Context(), claims.UserID)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":        "success",
			"user":        user,
			"message":     "Successfully returned user",
			"server_time": time.Now(),
		})
	}
}
