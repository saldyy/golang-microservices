package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/saldyy/golang-microservices/internal/health"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	defer stop()

	server := &Server{logger: slogger}

	go func() {
		if err := server.Run(":8080"); err != nil && err != http.ErrServerClosed {
			server.echo.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.echo.Shutdown(ctx); err != nil {
		server.logger.Error("Shutting down server")
	}
}

func (s *Server) Run(listen string) error {
	s.logger.Info("Configuring HTTP server")

	s.echo = echo.New()
	s.echo.HideBanner = true
	s.echo.Use(slogecho.New(s.logger))
	s.echo.Use(middleware.Recover())

  s.echo.GET("health", health.CheckHandler)

	return s.echo.Start(":8080")
}
