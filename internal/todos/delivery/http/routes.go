package http

import (
	"github.com/labstack/echo/v4"

	"github.com/Yangiboev/todo/internal/todos"
)

// Map todos routes
func MapToDosRoutes(todoGroup *echo.Group, h todos.Handlers) {
	todoGroup.POST("", h.Create())
	todoGroup.GET("/list", h.GetAll())
	todoGroup.DELETE("/:id", h.Delete())
	todoGroup.PUT("/:id", h.Update())
	todoGroup.GET("/:id", h.GetByID())
}
