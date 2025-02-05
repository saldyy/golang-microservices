package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
)

const ()

type Server struct {
	echo   *echo.Echo
	logger *slog.Logger
}

func main() {
	err := godotenv.Load(filepath.Join(".", ".env"))
	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err != nil {
		slogger.Error("Error loading .env file", err)
		panic(err)
	}

	server := &Server{logger: slogger}
	server.Run(":8080")
}

func (s *Server) Run(listen string) error {
	s.logger.Info("Configuring HTTP server")

	s.echo = echo.New()
	s.echo.HideBanner = true
	s.echo.Use(slogecho.New(s.logger))
	s.echo.Use(middleware.Recover())

	return s.echo.Start(":8080")
}
