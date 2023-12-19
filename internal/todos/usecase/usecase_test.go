package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Yangiboev/todo/internal/models"
	// "github.com/Yangiboev/todo/internal/todos/mock"
	"github.com/Yangiboev/todo/pkg/logger"
	"github.com/Yangiboev/todo/pkg/utils"
)

func TestToDosUC_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockToDoRepo := mock.NewMockRepository(ctrl)
	todoUC := NewToDosUseCase(nil, mockToDoRepo, apiLogger)

	comm := &models.ToDo{}

	mockToDoRepo.EXPECT().Create(context.Background(), gomock.Eq(comm)).Return(comm, nil)

	createdToDo, err := todoUC.Create(context.Background(), comm)
	require.NoError(t, err)
	require.NotNil(t, createdToDo)
}

func TestToDosUC_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockToDoRepo := mock.NewMockRepository(ctrl)
	todoUC := NewToDosUseCase(nil, mockToDoRepo, apiLogger)

	todoID := uuid.New()
	title := "title"

	todo := &models.ToDo{
		ToDoID: todoID,
		Title:  title,
	}

	baseToDo := &models.ToDo{
		Title: title,
	}

	mockToDoRepo.EXPECT().GetByID(context.Background(), gomock.Eq(todo.ToDoID)).Return(baseToDo, nil)
	mockToDoRepo.EXPECT().Update(context.Background(), gomock.Eq(todo)).Return(todo, nil)

	updatedToDo, err := todoUC.Update(context.Background(), todo)
	require.NoError(t, err)
	require.NotNil(t, updatedToDo)
}

func TestToDosUC_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockToDoRepo := mock.NewMockRepository(ctrl)
	todoUC := NewToDosUseCase(nil, mockToDoRepo, apiLogger)

	todoID := uuid.New()
	title := "title"

	todo := &models.ToDo{
		ToDoID: todoID,
		Title:  title,
	}

	baseToDo := &models.ToDo{
		Title: title,
	}

	mockToDoRepo.EXPECT().GetByID(context.Background(), gomock.Eq(todo.ToDoID)).Return(baseToDo, nil)
	mockToDoRepo.EXPECT().Delete(context.Background(), gomock.Eq(todo.ToDoID)).Return(nil)

	err := todoUC.Delete(context.TODO(), todo.ToDoID)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestToDosUC_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockToDoRepo := mock.NewMockRepository(ctrl)
	todoUC := NewToDosUseCase(nil, mockToDoRepo, apiLogger)

	todo := &models.ToDo{
		ToDoID: uuid.New(),
	}

	baseToDo := &models.ToDo{}

	ctx := context.Background()

	mockToDoRepo.EXPECT().GetByID(ctx, gomock.Eq(todo.ToDoID)).Return(baseToDo, nil)

	todoBase, err := todoUC.GetByID(ctx, todo.ToDoID)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, todoBase)
}

func TestToDosUC_GetAll(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockToDoRepo := mock.NewMockRepository(ctrl)
	todoUC := NewToDosUseCase(nil, mockToDoRepo, apiLogger)

	todoList := &models.ToDosList{}

	ctx := context.Background()
	query := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}

	mockToDoRepo.EXPECT().GetAll(ctx, query).Return(todoList, nil)

	todoLists, err := todoUC.GetAll(ctx, query)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, todoLists)
}
