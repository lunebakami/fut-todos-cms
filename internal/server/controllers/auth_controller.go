package controllers

import (
	"fut-todos-cms/internal/database"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/joho/godotenv/autoload"
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

type SignInResponse struct {
  Token string `json:"access_token"`
}

func NewAuthController(service database.Service) AuthController {
	return &auth_controller{
		service: service,
	}
}

func generateTokenJWT(sub string) (string, error) {
	claims := jwt.MapClaims{
		"sub": sub,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(2 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := os.Getenv("JWT_SECRET")
	secretKeyBytes := []byte(secretKey)
	tokenString, err := token.SignedString(secretKeyBytes)
	if err != nil {
		log.Printf("Error creating token: %e\n", err)
		return "", err
	}

	return tokenString, err
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

	token, err := generateTokenJWT(user.ID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, "Error creating token")
	}

	response := SignInResponse{
		Token: token,
	}

	return c.JSON(http.StatusOK, response)
}
