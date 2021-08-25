package httpd

import (
	"friend-management/internal/core/domain"
	"friend-management/internal/httpd/middleware/jwtauth"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type friendPutRequest struct {
	FriendID         string `json:"friend_id" binding:"required"`
	VersionRev       string `json:"version_rev" binding:"required"`
	FriendVersionRev string `json:"friend_version_rev" binding:"required"`
}

func (s *Server) FriendPut() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody friendPutRequest

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		claims, _ := c.Keys[jwtauth.JWTClaimsCtxKey].(jwtauth.JWTClaims)

		err := s.repo.AddFriend(c.Request.Context(), domain.AddFriendRequest{
			ID:               claims.UserID,
			FriendID:         requestBody.FriendID,
			VersionRev:       requestBody.VersionRev,
			FriendVersionRev: requestBody.FriendVersionRev,
		})
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":        "success",
			"message":     "Successfully add friend",
			"server_time": time.Now(),
		})
	}
}
