package util

import (
	"friend-management/internal/core/domain"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetPaginationQuery ...
func GetPaginationQuery(c *gin.Context) domain.UserQuery {
	const layout = "2006-01-02T15:04:05Z07:00"

	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	if page <= 0 {
		page = 1
	}
	perPage, _ := strconv.ParseInt(c.Query("per_page"), 10, 64)
	if perPage <= 0 {
		perPage = 200
	}
	rangeStartCreatedAt, _ := time.Parse(layout, c.Query("range_start_created_at"))
	rangeEndCreatedAt, _ := time.Parse(layout, c.Query("range_end_created_at"))

	return domain.UserQuery{
		Page:                page,
		PerPage:             perPage,
		Order:               c.DefaultQuery("order", "desc"),
		RangeStartCreatedAt: rangeStartCreatedAt,
		RangeEndCreatedAt:   rangeEndCreatedAt,
	}
}
