package handler

import (
  "net/http"
  "github.com/labstack/echo/v4"
)

func LoginHandler(c echo.Context) error {

  return c.String(http.StatusOK, "Signin")
}

func SignUpHandler(c echo.Context) error {

  return c.String(http.StatusOK, "Signup")
}
