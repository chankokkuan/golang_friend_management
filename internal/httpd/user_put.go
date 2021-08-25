package httpd

import (
	"friend-management/internal/core/domain"
	"friend-management/internal/httpd/middleware/jwtauth"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type userPutRequest struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	VersionRev string `json:"version_rev" binding:"required"`
}

func (s *Server) UserPut() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody := userPutRequest{}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		claims, _ := c.Keys[jwtauth.JWTClaimsCtxKey].(jwtauth.JWTClaims)

		req := domain.UpdateUserRequest{
			ID:         claims.UserID,
			Name:       requestBody.Name,
			Email:      requestBody.Email,
			VersionRev: requestBody.VersionRev,
		}

		user, err := s.repo.UpdateUser(c.Request.Context(), req)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":        "success",
			"user":        user,
			"message":     "Successfully updated user",
			"server_time": time.Now(),
		})
	}
}
