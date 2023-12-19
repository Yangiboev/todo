//go:generate mockgen -source pg_repository.go -destination mock/pg_repository_mock.go -package mock
package todos

import (
	"context"

	"github.com/google/uuid"

	"github.com/Yangiboev/todo/internal/models"
	"github.com/Yangiboev/todo/pkg/utils"
)

// ToDos repository interface
type Repository interface {
	Create(ctx context.Context, todo *models.ToDo) (*models.ToDo, error)
	Update(ctx context.Context, todo *models.ToDo) (*models.ToDo, error)
	Delete(ctx context.Context, todoID uuid.UUID) error
	GetByID(ctx context.Context, todoID uuid.UUID) (*models.ToDo, error)
	GetAll(ctx context.Context, title string, query *utils.PaginationQuery) (*models.ToDosList, error)
}
