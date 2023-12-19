package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/Yangiboev/todo/config"
	"github.com/Yangiboev/todo/internal/models"

	// "github.com/Yangiboev/todo/internal/todos/mock"
	"github.com/Yangiboev/todo/internal/todos/usecase"
	"github.com/Yangiboev/todo/pkg/converter"
	"github.com/Yangiboev/todo/pkg/logger"
)

func TestToDosHandlers_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cfg := config.Config{}
	apiLogger := logger.NewApiLogger(&cfg)
	mockToDoUC := mock.NewMockUseCase(ctrl)
	todoUC := usecase.NewToDosUseCase(&cfg, mockToDoUC, apiLogger)

	todoHandlers := NewToDosHandlers(&cfg, todoUC, apiLogger)
	handlerFunc := todoHandlers.Create()

	todoID := uuid.New()
	todo := &models.ToDo{
		ToDoID: todoID,
		Title:  "message Key: 'todo.Message' Error:Field validation for 'Message' failed on the 'gte' tag",
	}

	buf, err := converter.AnyToBytesBuffer(todo)
	require.NoError(t, err)
	require.NotNil(t, buf)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/v1/todos", strings.NewReader(buf.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(req, res)

	mockComm := &models.ToDo{
		ToDoID: todoID,
		Title:  "message",
	}

	mockToDoUC.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mockComm, nil)

	err = handlerFunc(ctx)
	require.NoError(t, err)
}

func TestToDosHandlers_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockToDoUC := mock.NewMockUseCase(ctrl)
	todoUC := usecase.NewToDosUseCase(nil, mockToDoUC, apiLogger)

	todoHandlers := NewToDosHandlers(nil, todoUC, apiLogger)
	handlerFunc := todoHandlers.GetByID()

	r := httptest.NewRequest(http.MethodGet, "/v1/todos/5c9a9d67-ad38-499c-9858-086bfdeaf7d2", nil)
	w := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(r, w)
	c.SetParamNames("id")
	c.SetParamValues("5c9a9d67-ad38-499c-9858-086bfdeaf7d2")

	comm := &models.ToDo{}

	mockToDoUC.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(comm, nil)

	err := handlerFunc(c)
	require.NoError(t, err)
}

func TestToDosHandlers_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cfg := config.Config{}

	apiLogger := logger.NewApiLogger(&cfg)
	mockToDoUC := mock.NewMockUseCase(ctrl)
	todoUC := usecase.NewToDosUseCase(&cfg, mockToDoUC, apiLogger)

	todoHandlers := NewToDosHandlers(&cfg, todoUC, apiLogger)
	handlerFunc := todoHandlers.Delete()

	todoID := uuid.New()
	commID := uuid.New()
	comm := &models.ToDo{
		ToDoID: todoID,
	}

	r := httptest.NewRequest(http.MethodDelete, "/v1/todos/5c9a9d67-ad38-499c-9858-086bfdeaf7d2", nil)
	w := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(r, w)
	c.SetParamNames("id")
	c.SetParamValues(commID.String())

	mockToDoUC.EXPECT().GetByID(gomock.Any(), commID).Return(comm, nil)
	mockToDoUC.EXPECT().Delete(gomock.Any(), commID).Return(nil)

	err := handlerFunc(c)
	require.NoError(t, err)
}
