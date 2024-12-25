package controllers

import (
	"fmt"
	"fut-todos-cms/internal/database"
	"fut-todos-cms/internal/server/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PostController interface {
	CreatePost(c echo.Context) error
	GetPosts(c echo.Context) error
}

type post_controller struct {
	service database.Service
}

type CreatePostResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
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

func (controller *post_controller) CreatePost(c echo.Context) error {
	p := models.Post{}

	if err := c.Bind(&p); err != nil {
		return err
	}

	result, err := controller.service.InsertPost(p)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Post created successfully. Rows affected: %d", rowsAffected))
}
