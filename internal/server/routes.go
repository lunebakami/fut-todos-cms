package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"fut-todos-cms/internal/server/controllers"

	_ "github.com/joho/godotenv/autoload"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	authController := controllers.NewAuthController(s.db)
	postController := controllers.NewPostController(s.db)
	userController := controllers.NewUserController(s.db)

  e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"*"},
  }))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", s.HelloWorldHandler)

	e.GET("/health", s.healthHandler)

	e.POST("/auth/signin", authController.SignIn)

  r := e.Group("")

  r.Use(AuthMiddleware)

	r.GET("/posts", postController.GetPosts)
	r.POST("/posts", postController.CreatePost)
  r.DELETE("/posts/:id", postController.DeletePost)

	r.GET("/users", userController.GetUsers)
	r.POST("/users", userController.CreateUser)
  r.DELETE("/users/:id", userController.DeleteUser)

	return e
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "Missing or invalid Authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "Invalid token format.",
			})
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unexpected signing method")
			}

			secretKey := os.Getenv("JWT_SECRET")
      secretKeyBytes := []byte(secretKey)
			return secretKeyBytes, nil
		})

    if err != nil {
      return c.JSON(http.StatusUnauthorized, ErrorResponse{
        Error: "Invalid or expired token",
      })
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
      c.Set("user", claims)
    } else {
      return c.JSON(http.StatusUnauthorized, ErrorResponse{
        Error: "Invalid token claims",
      })
    }

    return next(c)
	}
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
