package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/saldyy/golang-microservices/config/database"
	"github.com/saldyy/golang-microservices/handler"
	slogecho "github.com/samber/slog-echo"
)

type Server struct {
	echo   *echo.Echo
	logger *slog.Logger
}

func main() {

  fmt.Printf("%s\n", filepath.Join(".", ".env"))
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
	e := echo.New()
	e.HideBanner = true
	e.Use(slogecho.New(s.logger))
	e.Use(middleware.Recover())

	e.GET("/health", handler.HealthCheckHandler)

	e.POST("login", handler.LoginHandler)
	e.POST("/sign-up", handler.SignUpHandler)

	e.GET("/users", handler.GetUser)

	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		code := 500
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		s.logger.Warn("HTTP request error", "statusCode", code, "path", ctx.Path(), "err", err)
		ctx.Response().WriteHeader(code)
	}

	s.echo = e

	database.Instance = database.InitMongoClient()

	return s.echo.Start(":8080")
}
