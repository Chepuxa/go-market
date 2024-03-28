package repositories

import (
	"database/sql"
	"log"
	"training/proj/internal/api/models"
)

type CategoryRepositoryInterface interface {
	Create(*models.Category) (models.Category, error)
	GetAll() ([]models.Category, error)
	GetByName(string) (models.Category, error)
	GetById(int64) (models.Category, error)
	Delete(int64) ([]models.Category, error)
	Update(int64, *models.Category) (int64, error)
	GetCategoryItems(int64) ([]models.Item, error)
}

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) Create(categoryReq *models.Category) (models.Category, error) {
	sqlStatement := `INSERT INTO categories (category) VALUES ($1) RETURNING *`

	var categoryResp models.Category

	err := r.db.QueryRow(sqlStatement, categoryReq.Category).Scan(&categoryResp.CategoryID, &categoryResp.Category)

	return categoryResp, err
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	categories := make([]models.Category, 0)

	sqlStatement := `SELECT * FROM categories`

	rows, queryErr := r.db.Query(sqlStatement)

	if queryErr != nil {
		return nil, queryErr
	}

	defer rows.Close()

	for rows.Next() {
		var category models.Category

		scanErr := rows.Scan(&category.CategoryID, &category.Category)

		if scanErr != nil {
			return nil, scanErr
		}

		categories = append(categories, category)

	}

	return categories, nil
}

func (r *CategoryRepository) Delete(id int64) (int64, error) {
	sqlStatement := `DELETE FROM categories WHERE category_id = $1`

	res, execErr := r.db.Exec(sqlStatement, id)

	if execErr != nil {
		return 0, execErr
	}

	rowsAffected, rowsErr := res.RowsAffected()

	if rowsErr != nil {
		log.Fatalf("Daaaym what's that")
	}

	return rowsAffected, nil
}

func (r *CategoryRepository) Update(id int64, categoryReq *models.Category) (models.Category, error) {
	sqlStatement := `UPDATE categories SET category = $2 WHERE category_id = $1 RETURNING *`

	var categoryResp models.Category

	err := r.db.QueryRow(sqlStatement, id, categoryReq.Category).Scan(&categoryResp.CategoryID, &categoryResp.Category)

	return categoryResp, err
}

func (r *CategoryRepository) GetById(id int64) (models.Category, error) {
	var category models.Category

	getCategoryStatement := `SELECT * FROM categories WHERE category_id = $1`

	row := r.db.QueryRow(getCategoryStatement, id)

	err := row.Scan(&category.CategoryID, &category.Category)

	return category, err
}

func (r *CategoryRepository) GetByName(name string) (models.Category, error) {
	var category models.Category

	getCategoryStatement := `SELECT * FROM categories WHERE category = $1`

	row := r.db.QueryRow(getCategoryStatement, name)

	err := row.Scan(&category.CategoryID, &category.Category)

	return category, err
}

func (r *CategoryRepository) GetCategoryItems(id int64) ([]models.Item, error) {
	items := make([]models.Item, 0)

	_, getErr := r.GetById(id)

	if getErr != nil {
		return nil, getErr
	}

	sqlStatement := `SELECT item_id, item, price FROM items
	INNER JOIN categories_items
	USING (item_id)
	WHERE category_id = $1`

	rows, queryErr := r.db.Query(sqlStatement, id)

	if queryErr != nil {
		return nil, queryErr
	}

	defer rows.Close()

	for rows.Next() {
		var item models.Item

		scanErr := rows.Scan(&item.ItemID, &item.Item, &item.Price)

		if scanErr != nil {
			return nil, scanErr
		}

		items = append(items, item)

	}

	return items, nil
}
