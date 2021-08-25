package jwtauth

import (
	"strings"

	"friend-management/internal/core/domain"
	"friend-management/internal/log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTClaims ...
type JWTClaims struct {
	UserID string   `json:"user_id"`
	Name   string   `json:"name"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`

	jwt.StandardClaims
}

const (
	NoToken         = "jwt-no-token"
	BadToken        = "jwt-bad-token"
	ClaimsNotFound  = "no-jwtclaims-in-context"
	ClaimsIncorrect = "jwtclaims-wrong-type"
)

var (
	ErrNoToken         = domain.CustomizeError(NoToken, "Bearer token is not found in HTTP Authorization header", nil)
	ErrBadToken        = domain.CustomizeError(BadToken, "Invalid JWT token or signature", nil)
	ErrClaimsNotFound  = domain.CustomizeError(ClaimsNotFound, "JWT Claims not found in context", nil)
	ErrClaimsIncorrect = domain.CustomizeError(ClaimsIncorrect, "JWT Claims have incorrect data structure", nil)
)

// JWTClaimsCtxKey ...
const JWTClaimsCtxKey string = "JWTClaims"

// JWTMiddleware handle JWT Auth process
func JWTMiddleware(jwtKey []byte, logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.Request.Header.Get("Authorization")
		splitToken := strings.Split(tokenString, "Bearer ")

		if len(splitToken) <= 1 {
			// no bearer token found
			c.Error(ErrNoToken)

			c.Abort()
			return
		}

		tokenString = splitToken[1]

		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// function to get the jwt signing key
			return jwtKey, nil
		})

		if err != nil {
			c.Error(ErrBadToken)

			c.Abort()
			return
		}

		claims := token.Claims.(*JWTClaims)
		newCtx := logger.Attach(c.Request.Context(), "user_id", claims.UserID)
		c.Request = c.Request.WithContext(newCtx)

		if c.Keys == nil {
			c.Keys = make(map[string]interface{})
		}
		c.Keys[JWTClaimsCtxKey] = *claims

		c.Next()
	}
}
