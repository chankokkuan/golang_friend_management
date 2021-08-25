package httpd

import (
	"friend-management/internal/httpd/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) UsersGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := util.GetPaginationQuery(c)
		results, err := s.repo.GetUsers(c.Request.Context(), query)

		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":  "success",
			"meta":  results.Meta,
			"users": results.Users,
		})
	}
}
