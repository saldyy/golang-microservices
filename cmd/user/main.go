package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	//	"path/filepath"
	"time"

	//	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/saldyy/golang-microservices/internal/database"
	"github.com/saldyy/golang-microservices/internal/health"
	slogecho "github.com/samber/slog-echo"
)

type Server struct {
	echo   *echo.Echo
	logger *slog.Logger
	db     database.Db
}

func main() {
	// err := godotenv.Load(filepath.Join(".", ".env"))
	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// if err != nil {
	// 	slogger.Error("Error loading .env file", err)
	// 	panic(err)
	// }

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	defer stop()

	db := database.New(os.Getenv("PG_CONNECTION_STRING"))
	server := &Server{logger: slogger, db: *db}

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

	s.RegisterRoutes()

	return s.echo.Start(":8080")
}

func (s *Server) RegisterRoutes() {

	healthHandler := health.New(&s.db)

	s.echo.GET("health", healthHandler.CheckHandler)
}
