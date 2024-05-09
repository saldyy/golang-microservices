package user

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/saldyy/golang-microservices/model"
	"github.com/saldyy/golang-microservices/repository"
	"github.com/saldyy/golang-microservices/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userRepo         *repository.UserRepository
	jwtBlacklistRepo *repository.JwtBlackListRepository
	logger           *slog.Logger
}

type LoginCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *AuthHandler) LoginHandler(c echo.Context) error {
	credential := new(LoginCredential)

	if err := c.Bind(credential); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid credential"})
	}

	user, err := u.userRepo.FindByUsername(credential.Username)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			u.logger.Debug("User not found", err)
			return c.JSON(http.StatusForbidden, map[string]string{"message": "User not found"})
		}

		return c.String(http.StatusOK, "Error")
	}

	isPasswordMatch := utils.IsMatchBcryptHash(user.Password, credential.Password)

	if !isPasswordMatch {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "Invalid password"})
	}

	jwtClaims := &jwt.RegisteredClaims{
		Subject:   user.Id,
		Issuer:    "auth-service",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtClaims)
	token.Header["kid"] = os.Getenv("JWT_KID")

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		u.logger.Error("Error signing token", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"access_token": t})
}

func (u *AuthHandler) RegisterHandler(c echo.Context) error {
	credential := new(LoginCredential)

	if err := c.Bind(credential); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid payload"})
	}

	existedUser, err := u.userRepo.FindByUsername(credential.Username)

	if existedUser != nil {
		u.logger.Error("User already existed", err)
		return c.JSON(http.StatusConflict, map[string]string{"message": "User already existed"})
	}

	hashedPassword, err := utils.BcryptHash(credential.Password)

	if err != nil {
		u.logger.Error("Error hashing password", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	user := &model.User{Username: credential.Username, Password: hashedPassword}
	_, err = u.userRepo.InsertUser(user)

	if err != nil {
		u.logger.Error("Error inserting user", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return c.String(http.StatusOK, "user.Username")
}

func (u *AuthHandler) LogoutHandler(c echo.Context) error {
	auth := c.Request().Header.Get("Authorization")
	parts := strings.Split(auth, " ")
	jwtToken := parts[1]

	if jwtToken == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid token"})
	}

	err := u.jwtBlacklistRepo.Add(jwtToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return c.NoContent(http.StatusOK)
}

func RegisterAuthHandlers(e *echo.Echo, repo *repository.Repository, logger *slog.Logger) *AuthHandler {
	authHandler := AuthHandler{userRepo: repo.UserRepository, jwtBlacklistRepo: repo.JwtBlackListRepository, logger: logger}

	e.POST("/login", authHandler.LoginHandler)
	e.POST("/register", authHandler.RegisterHandler)
	e.POST("/logout", authHandler.LogoutHandler)

	return &authHandler
}
