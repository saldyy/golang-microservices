package main

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/saldyy/golang-microservices/config/database"
	"github.com/saldyy/golang-microservices/handler"
	"github.com/saldyy/golang-microservices/repository"
	slogecho "github.com/samber/slog-echo"
)

const ()

type Server struct {
	echo   *echo.Echo
	logger *slog.Logger
}

func main() {
	err := godotenv.Load(filepath.Join(".", ".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	server := &Server{logger: slogger}
	server.Run(":8080")
}

func (s *Server) Run(listen string) error {
	keyId := os.Getenv("JWT_KID")
	jwtSecret := os.Getenv("JWT_SECRET")

	s.logger.Info("Configuring HTTP server")

	database.Instance = database.InitMongoClient(s.logger)
	database.RedisInstance = database.NewRedisClientInstance(s.logger)
	repos := repository.Init(database.Instance.DB, database.RedisInstance.Client)


	e := echo.New()
	e.HideBanner = true
	e.Use(slogecho.New(s.logger))
	e.Use(middleware.Recover())

  // Check token in blacklist
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/login" || c.Path() == "/register" || c.Path() == "/health" {
				return next(c)
			}
			auth := c.Request().Header.Get("Authorization")
			parts := strings.Split(auth, " ")
			jwtToken := parts[1]
			blacklistToken, err := repos.JwtBlackListRepository.Get(jwtToken)

			if err != nil {
				s.logger.Error("Error getting token from blacklist", err)
			}
			if blacklistToken == "1" {
				return c.JSON(401, map[string]string{"message": "Unauthorized"})
			}
			return next(c)
		}
	})

	e.Use(echojwt.WithConfig(echojwt.Config{
		ContextKey:  "user",
		SigningKeys: map[string]interface{}{(keyId): []byte(jwtSecret)},
		Skipper: func(c echo.Context) bool {
			if c.Path() == "/login" || c.Path() == "/register" || c.Path() == "/health" {
				return true
			}
			return false
		},
		ErrorHandler: func(c echo.Context, er error) error {
      s.logger.Error("Error validating token", er)
			return c.JSON(401, map[string]string{"message": "Unauthorized"})
		},
		TokenLookup:   "header:Authorization:Bearer ",
		SigningMethod: jwt.SigningMethodHS512.Name,
	}))

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
