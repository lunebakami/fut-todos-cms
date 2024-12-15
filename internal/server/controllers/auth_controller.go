package controllers

import (
	"fut-todos-cms/internal/database"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthController interface {
	SignIn(c echo.Context) error
}

type auth_controller struct {
	service database.Service
}

type SignInDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthController(service database.Service) AuthController {
	return &auth_controller{
		service: service,
	}
}

func (controller *auth_controller) SignIn(c echo.Context) error {
	u := new(SignInDTO)
	if err := c.Bind(u); err != nil {
		return err
	}

	user, err := controller.service.GetUserByEmail(u.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid email")
	}

	inputPassword := []byte(u.Password)
	userPassword := []byte(user.Password)

	err = bcrypt.CompareHashAndPassword(
		userPassword,
		inputPassword,
	)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid password")
	}

	return c.JSON(http.StatusOK, user)
}
