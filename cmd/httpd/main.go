package main

import (
	"fmt"
	"friend-management/internal/config"
	"friend-management/internal/httpd"
	"friend-management/internal/log"
	"friend-management/internal/repository"
	"friend-management/pkg/mongo"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Startup error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		cfg           = config.LoadConfig()
		jwtSigningKey = []byte(cfg.JWTKey)
	)

	logger := log.NewGokitWrapper("friend-management")

	// init datadog tracer
	tracer.Start(
		tracer.WithServiceName("friend-management"),
		tracer.WithGlobalTag("env", cfg.Environment),
	)
	defer tracer.Stop()

	conn, cli, err := mongo.NewConnection(
		cfg.MongoDB.Host,
		cfg.MongoDB.Port,
		cfg.MongoDB.Database,
		cfg.MongoDB.Username,
		cfg.MongoDB.Password,
	)
	if err != nil {
		logger.Error("msg", "Error connecting to MongoDB", "error", err)
		return err
	}
	logger.Info("msg", fmt.Sprintf("Connected to MongoDB Database: %s", cfg.MongoDB.Database))

	repo := repository.NewUserRepo(conn, cli, logger)

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Same format as default gin.Logger(), with TraceID injected into the logs
		span, _ := tracer.SpanFromContext(param.Request.Context())
		return fmt.Sprintf("[GIN] %v | %3d | %13vms | %15s | %-7s %#v | %v | %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.StatusCode,
			float64(param.Latency)/float64(time.Millisecond),
			param.ClientIP,
			param.Method,
			param.Path,
			span.Context().TraceID(),
			param.ErrorMessage,
		)
	}))

	router.Use(gintrace.Middleware("friend-management"))
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("x-correlation-id")
	corsConfig.AddAllowHeaders("authorization")
	router.Use(cors.New(corsConfig))

	server := httpd.NewServer(
		router, repo, jwtSigningKey, logger,
	)

	err = server.Run()
	if err != nil {
		return err
	}
	return nil
}
