package controllers

import (
	"fut-todos-cms/internal/database"
	"fut-todos-cms/internal/server/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserController interface {
	CreateUser(c echo.Context) error
	GetUsers(c echo.Context) error
}

type user_controller struct {
	service database.Service
}

func NewUserController(service database.Service) UserController {
	return &user_controller{
		service: service,
	}
}

func (controller *user_controller) GetUsers(c echo.Context) error {
  users, err := controller.service.GetUsers()

  if err != nil {
    log.Println(err)
    return err
  }

  return c.JSON(http.StatusOK, users)
}

func (controller *user_controller) CreateUser(c echo.Context) error {
  u := models.User{}

  if err := c.Bind(&u); err != nil {
    return err
  }

  password := []byte(u.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

  u.Password = string(hashedPassword)

  result, err := controller.service.InsertUser(u)
  if err != nil {
    return nil
  }

  return c.JSON(http.StatusOK, result)
}


