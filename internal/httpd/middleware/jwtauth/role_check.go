package jwtauth

import (
	"friend-management/internal/core/domain"

	"github.com/gin-gonic/gin"
)

// RoleCheckMiddleware ...
func RoleCheckMiddleware(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {

		v, ok := c.Keys[JWTClaimsCtxKey]

		if !ok {
			// this should not happen.
			c.Error(ErrClaimsNotFound)
			c.Abort()
			return
		}

		claims, ok := v.(JWTClaims)
		if !ok {
			// this should not happen.
			c.Error(ErrClaimsIncorrect)
			c.Abort()
			return
		}

		// TODO: not a very efficient of doing this.
		m := make(map[string]bool)
		for _, item := range allowedRoles {
			m[item] = true
		}
		for _, item := range claims.Roles {
			if _, ok := m[item]; ok {
				// matching role found.
				c.Next()
				return
			}
		}

		// User access denied.
		c.Error(domain.ErrAccessDenied)
		c.Abort()
	}
}
