package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"fut-todos-cms/internal/server/controllers"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

  authController := controllers.NewAuthController(s.db)
  postController := controllers.NewPostController(s.db)
  userController := controllers.NewUserController(s.db)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", s.HelloWorldHandler)

	e.GET("/health", s.healthHandler)

	e.GET("/posts", postController.GetPosts)
  e.POST("/posts", postController.CreatePost)

  e.GET("/users", userController.GetUsers)
  e.POST("/users", userController.CreateUser)

  e.POST("/auth/signin", authController.SignIn)
	return e
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
