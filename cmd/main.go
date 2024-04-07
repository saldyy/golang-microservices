package main

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/saldyy/golang-microservices/config/database"
	"github.com/saldyy/golang-microservices/handler"
	"github.com/saldyy/golang-microservices/repository"
	slogecho "github.com/samber/slog-echo"
)

type Server struct {
	echo   *echo.Echo
	logger *slog.Logger
}

func main() {
	err := godotenv.Load(filepath.Join(".", ".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	slogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	server := &Server{logger: slogger}
	server.Run(":8080")
}

func (s *Server) Run(listen string) error {
	s.logger.Info("Configuring HTTP server")

	database.Instance = database.InitMongoClient()
  repos := repository.Init(database.Instance.DB)

	e := echo.New()
	e.HideBanner = true
	e.Use(slogecho.New(s.logger))
	e.Use(middleware.Recover())


	e.GET("/health", handler.HealthCheckHandler)
  handler.RegisterRoutes(e, repos, s.logger)



	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		code := 500
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		s.logger.Warn("HTTP request error", "statusCode", code, "path", ctx.Path(), "err", err)
		ctx.Response().WriteHeader(code)
	}

	s.echo = e


	return s.echo.Start(":8080")
}
