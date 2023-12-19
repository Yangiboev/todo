package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/Yangiboev/todo/internal/models"
	"github.com/Yangiboev/todo/internal/todos"
	"github.com/Yangiboev/todo/pkg/utils"
)

// ToDos Repository
type todosRepo struct {
	db *sqlx.DB
}

// ToDos Repository constructor
func NewToDosRepository(db *sqlx.DB) todos.Repository {
	return &todosRepo{db: db}
}

// Create todo
func (r *todosRepo) Create(ctx context.Context, todo *models.ToDo) (*models.ToDo, error) {
	newUUID := uuid.New()
	c := &models.ToDo{}
	createToDo := `INSERT INTO todos (id, title) VALUES ($1, $2) RETURNING *`
	if err := r.db.QueryRowxContext(
		ctx,
		createToDo,
		newUUID,
		&todo.Title,
	).StructScan(c); err != nil {
		return nil, errors.Wrap(err, "todosRepo.Create.StructScan")
	}

	return c, nil
}

// Update todo
func (r *todosRepo) Update(ctx context.Context, todo *models.ToDo) (*models.ToDo, error) {
	updateToDo := `UPDATE todos SET title = $1 WHERE id = $2 RETURNING *`
	res := &models.ToDo{}
	if err := r.db.QueryRowxContext(ctx, updateToDo, todo.Title, todo.ToDoID).StructScan(res); err != nil {
		return nil, errors.Wrap(err, "todosRepo.Update.QueryRowxContext")
	}

	return res, nil
}

// Delete todo
func (r *todosRepo) Delete(ctx context.Context, toDoID uuid.UUID) error {
	deleteToDo := `DELETE FROM todos WHERE id = $1`

	result, err := r.db.ExecContext(ctx, deleteToDo, toDoID)
	if err != nil {
		return errors.Wrap(err, "todosRepo.Delete.ExecContext")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "todosRepo.Delete.RowsAffected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "todosRepo.Delete.rowsAffected")
	}

	return nil
}

// GetByID todo
func (r *todosRepo) GetByID(ctx context.Context, toDoID uuid.UUID) (*models.ToDo, error) {
	getToDoByID := `SELECT id, title, created_at
	FROM todos 
	WHERE id = $1`
	todo := &models.ToDo{}
	if err := r.db.GetContext(ctx, todo, getToDoByID, toDoID); err != nil {
		return nil, errors.Wrap(err, "todosRepo.GetByID.GetContext")
	}
	return todo, nil
}

// GetAll ToDos
func (r *todosRepo) GetAll(ctx context.Context, title string, query *utils.PaginationQuery) (*models.ToDosList, error) {
	var (
		totalCount    int
		getTotalCount = `SELECT COUNT(id) FROM todos WHERE 1=1`
		getAllToDos   = `SELECT id, title, created_at
							FROM todos where 1=1`
	)
	if title != "" {
		getTotalCount = fmt.Sprintf("%s%s", getTotalCount, " and title LIKE '%"+title+"%';")
		getAllToDos = fmt.Sprintf("%s%s", getAllToDos, " and title LIKE '%"+title+"%' ")
	}
	getAllToDos += " ORDER BY created_at OFFSET $1 LIMIT $2;"
	if err := r.db.QueryRowContext(ctx, getTotalCount).Scan(&totalCount); err != nil {
		return nil, errors.Wrap(err, "todosRepo.GetAll.QueryRowContext")
	}

	if totalCount == 0 {
		return &models.ToDosList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
			Page:       query.GetPage(),
			Size:       query.GetSize(),
			HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			ToDos:      make([]*models.ToDo, 0),
		}, nil
	}

	rows, err := r.db.QueryxContext(ctx, getAllToDos, query.GetOffset(), query.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "todosRepo.GetAll.QueryxContext")
	}
	defer rows.Close()

	todosList := make([]*models.ToDo, 0, query.GetSize())
	for rows.Next() {
		todo := &models.ToDo{}
		if err = rows.StructScan(todo); err != nil {
			return nil, errors.Wrap(err, "todosRepo.GetAll.StructScan")
		}
		todosList = append(todosList, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "todosRepo.GetAll.rows.Err")
	}

	return &models.ToDosList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
		Page:       query.GetPage(),
		Size:       query.GetSize(),
		HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		ToDos:      todosList,
	}, nil
}
