package httpd

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LivenessGet ...
func (s *Server) LivenessGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		isConnectedToDB := s.repo.ConnectionCheck(c)

		if !isConnectedToDB {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"is_alive": false,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"is_alive": true,
		})
	}
}
