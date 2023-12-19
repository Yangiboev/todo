package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/Yangiboev/todo/internal/models"
)

func TestToDosRepo_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	commRepo := NewToDosRepository(sqlxDB)

	t.Run("Create", func(t *testing.T) {
		todoID := uuid.New()
		title := "title"

		rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(todoID, title)

		todo := &models.ToDo{
			ToDoID: todoID,
			Title:  title,
		}

		mock.ExpectQuery(createToDo).WithArgs(todo.ToDoID, todo.Title).WillReturnRows(rows)

		createdToDo, err := commRepo.Create(context.Background(), todo)

		require.NoError(t, err)
		require.NotNil(t, createdToDo)
		require.Equal(t, createdToDo, todo)
	})

	t.Run("Create ERR", func(t *testing.T) {
		newsUID := uuid.New()
		title := "title"
		createErr := errors.New("Create todo error")

		todo := &models.ToDo{
			ToDoID: newsUID,
			Title:  title,
		}

		mock.ExpectQuery(createToDo).WithArgs(todo.ToDoID, todo.Title).WillReturnError(createErr)

		createdToDo, err := commRepo.Create(context.Background(), todo)

		require.Nil(t, createdToDo)
		require.NotNil(t, err)
	})
}

func TestToDosRepo_Update(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	commRepo := NewToDosRepository(sqlxDB)

	t.Run("Update", func(t *testing.T) {
		todoID := uuid.New()
		title := "title"

		rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(todoID, title)

		todo := &models.ToDo{
			ToDoID: todoID,
			Title:  title,
		}

		updatedToDo, err := commRepo.Update(context.Background(), todo)
		mock.ExpectQuery(updateToDo).WithArgs(todo.ToDoID, todo.Title).WillReturnRows(rows)

		require.NoError(t, err)
		require.NotNil(t, updatedToDo)
		require.Equal(t, updatedToDo.ToDoID, todo.ToDoID)
	})

	t.Run("Update ERR", func(t *testing.T) {
		todoID := uuid.New()
		title := "title"

		todo := &models.ToDo{
			ToDoID: todoID,
			Title:  title,
		}

		updatedToDo, err := commRepo.Update(context.Background(), todo)

		require.NotNil(t, err)
		require.Nil(t, updatedToDo)
	})
}

func TestToDosRepo_Delete(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	commRepo := NewToDosRepository(sqlxDB)

	t.Run("Delete", func(t *testing.T) {
		todoID := uuid.New()
		mock.ExpectExec(deleteToDo).WithArgs(todoID).WillReturnResult(sqlmock.NewResult(1, 1))
		err := commRepo.Delete(context.Background(), todoID)

		require.NoError(t, err)
	})

	t.Run("Delete Err", func(t *testing.T) {
		todoID := uuid.New()

		mock.ExpectExec(deleteToDo).WithArgs(todoID).WillReturnResult(sqlmock.NewResult(1, 0))

		err := commRepo.Delete(context.Background(), todoID)
		require.NotNil(t, err)
	})
}
