package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/Yangiboev/todo/config"
	"github.com/Yangiboev/todo/internal/models"
	"github.com/Yangiboev/todo/internal/todos"
	"github.com/Yangiboev/todo/pkg/logger"
	"github.com/Yangiboev/todo/pkg/utils"
)

// ToDos UseCase
type todosUC struct {
	cfg       *config.Config
	todosRepo todos.Repository
	logger    logger.Logger
}

// ToDos UseCase constructor
func NewToDosUseCase(cfg *config.Config, todosRepo todos.Repository, logger logger.Logger) todos.UseCase {
	return &todosUC{cfg: cfg, todosRepo: todosRepo, logger: logger}
}

// Create todo
func (u *todosUC) Create(ctx context.Context, todo *models.ToDo) (*models.ToDo, error) {
	return u.todosRepo.Create(ctx, todo)
}

// Update todo
func (u *todosUC) Update(ctx context.Context, todo *models.ToDo) (*models.ToDo, error) {
	updatedToDo, err := u.todosRepo.Update(ctx, todo)
	if err != nil {
		return nil, err
	}

	return updatedToDo, nil
}

// Delete todo
func (u *todosUC) Delete(ctx context.Context, todoID uuid.UUID) error {

	if err := u.todosRepo.Delete(ctx, todoID); err != nil {
		return err
	}

	return nil
}

// GetByID todo
func (u *todosUC) GetByID(ctx context.Context, todoID uuid.UUID) (*models.ToDo, error) {

	return u.todosRepo.GetByID(ctx, todoID)
}

// GetAll todos
func (u *todosUC) GetAll(ctx context.Context, title string, query *utils.PaginationQuery) (*models.ToDosList, error) {
	return u.todosRepo.GetAll(ctx, title, query)
}
