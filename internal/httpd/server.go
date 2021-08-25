package httpd

import (
	"friend-management/internal/core/port"
	"friend-management/internal/httpd/middleware"
	"friend-management/internal/httpd/middleware/jwtauth"
	"friend-management/internal/log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router        *gin.Engine
	repo          port.UserRepository
	jwtSigningKey []byte
	logger        log.Logger
}

func NewServer(
	router *gin.Engine,
	repo port.UserRepository,
	jwtSigningKey []byte,
	logger log.Logger,
) *Server {
	return &Server{
		router,
		repo,
		jwtSigningKey,
		logger,
	}
}

var UserRole = []string{"User"}

func (s *Server) Run() error {
	router := s.router

	router.Use(middleware.LogRequestAndResponse(s.logger))
	router.Use(middleware.RespondWithError(s.logger))

	router.GET("/liveness", s.LivenessGet())
	router.POST("/user", s.UserPost())

	loginRequired := router.Group(".")
	loginRequired.Use(jwtauth.JWTMiddleware(s.jwtSigningKey, s.logger))
	{
		userOnlyPath := loginRequired.Group(".")
		userOnlyPath.Use(jwtauth.RoleCheckMiddleware(UserRole))

		userOnlyPath.GET("/user", s.UsersGet())
		userOnlyPath.GET("/user/:id", s.UserGet())
		userOnlyPath.PUT("/user", s.UserPut())
		userOnlyPath.PUT("/user/friend", s.FriendPut())
	}

	return router.Run()
}
