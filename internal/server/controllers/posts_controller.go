package controllers

import (
	"fmt"
	"fut-todos-cms/internal/database"
	"fut-todos-cms/internal/server/models"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type PostController interface {
	CreatePost(c echo.Context) error
	GetPosts(c echo.Context) error
	DeletePost(c echo.Context) error
}

type post_controller struct {
	service database.Service
}

func NewPostController(service database.Service) PostController {
	return &post_controller{
		service: service,
	}
}

func (controller *post_controller) GetPosts(c echo.Context) error {
	posts, err := controller.service.GetPosts()

	if err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, posts)
}

type CreatePostResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}

func (controller *post_controller) CreatePost(c echo.Context) error {
	p := models.Post{}
	user := c.Get("user")

	if err := c.Bind(&p); err != nil {
		return err
	}

	p.AuthorID = user.(jwt.MapClaims)["sub"].(string)

	result, err := controller.service.InsertPost(p)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error creating post: %s", err))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error getting rows affected: %s", err))
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Post created successfully. Rows affected: %d", rowsAffected))
}



func (controller *post_controller) DeletePost(c echo.Context) error {
  ID := c.Param("id")

  result, err := controller.service.DeletePost(ID)
  if err != nil {
    return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error deleting post: %s", err))
  }

  return c.JSON(http.StatusOK, fmt.Sprintf("Post deleted successfully. Rows affected: %d", result))
}
