package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/Yangiboev/todo/config"
	"github.com/Yangiboev/todo/internal/models"
	"github.com/Yangiboev/todo/internal/todos"
	"github.com/Yangiboev/todo/pkg/httpErrors"
	"github.com/Yangiboev/todo/pkg/logger"
	"github.com/Yangiboev/todo/pkg/utils"
)

// ToDos handlers
type todosHandlers struct {
	cfg     *config.Config
	todosUC todos.UseCase
	logger  logger.Logger
}

// NewToDosHandlers ToDos handlers constructor
func NewToDosHandlers(cfg *config.Config, todosUC todos.UseCase, logger logger.Logger) todos.Handlers {
	return &todosHandlers{cfg: cfg, todosUC: todosUC, logger: logger}
}

// Create
// @Summary Create new todo
// @Description create new todo
// @Tags ToDos
// @Accept  json
// @Produce  json
// @Param body body models.ToDoSwagger true "body"
// @Success 201 {object} models.ToDo
// @Failure 500 {object} httpErrors.RestErr
// @Router /todos [post]
func (h *todosHandlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {

		todo := &models.ToDo{}

		if err := utils.SanitizeRequest(c, todo); err != nil {
			return utils.ErrResponseWithLog(c, h.logger, err)
		}

		createdToDo, err := h.todosUC.Create(c.Request().Context(), todo)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, createdToDo)
	}
}

// Update
// @Summary Update todo
// @Description update new todo
// @Tags ToDos
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Param body body models.ToDoSwagger true "body"
// @Success 200 {object} models.ToDoSwagger
// @Failure 500 {object} httpErrors.RestErr
// @Router /todos/{id} [put]
func (h *todosHandlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {

		todosID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		comm := &models.ToDo{}
		if err = utils.SanitizeRequest(c, comm); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		updatedToDo, err := h.todosUC.Update(c.Request().Context(), &models.ToDo{
			ToDoID: todosID,
			Title:  comm.Title,
		})
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedToDo)
	}
}

// Delete
// @Summary Delete todo
// @Description delete todo
// @Tags ToDos
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {string} string	"ok"
// @Failure 500 {object} httpErrors.RestErr
// @Router /todos/{id} [delete]
func (h *todosHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {

		todosID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err = h.todosUC.Delete(c.Request().Context(), todosID); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// GetByID
// @Summary Get todo
// @Description Get todo by id
// @Tags ToDos
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} models.ToDo
// @Failure 500 {object} httpErrors.RestErr
// @Router /todos/{id} [get]
func (h *todosHandlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {

		todosID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		todo, err := h.todosUC.GetByID(c.Request().Context(), todosID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, todo)
	}
}

// GetAll
// @Summary Get ToDos
// @Description Get all todo
// @Tags ToDos
// @Accept  json
// @Produce  json
// @Param title query string false "title"
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Success 200 {object} models.ToDosList
// @Failure 500 {object} httpErrors.RestErr
// @Router /todos/list [get]
func (h *todosHandlers) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		toDoList, err := h.todosUC.GetAll(c.Request().Context(), c.QueryParam("title"), pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, toDoList)
	}
}
