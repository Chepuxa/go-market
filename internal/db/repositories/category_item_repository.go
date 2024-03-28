package repositories

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type CategoryItemrepositoryInterface interface {
	Create(int64, int64) *pgconn.PgError
}

type CategoryItemRepository struct {
	db *sql.DB
}

func NewCategoryItemRepository(db *sql.DB) *CategoryItemRepository {
	return &CategoryItemRepository{
		db: db,
	}
}

func (r *CategoryItemRepository) Create(categoryId int64, itemId int64) *pgconn.PgError {
	sqlStatement := `INSERT INTO categories_items (category_id, item_id) VALUES ($1, $2)`

	_, err := r.db.Exec(sqlStatement, categoryId, itemId)

	var e *pgconn.PgError
	errors.As(err, &e)

	return e
}
