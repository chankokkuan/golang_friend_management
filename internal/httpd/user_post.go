package httpd

import (
	"friend-management/internal/core/domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type userPostRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required"`
}

func (s *Server) UserPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody := userPostRequest{}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, err := s.repo.CreateUser(c.Request.Context(), domain.CreateUserRequest{
			Name:  requestBody.Name,
			Email: requestBody.Email,
		})
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"code":        "success",
			"message":     "Successfully created user",
			"user":        user,
			"server_time": time.Now(),
		})
	}
}
