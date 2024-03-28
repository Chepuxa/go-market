package repositories

import (
	"database/sql"
	"log"
	"training/proj/internal/api/models"
)

type ItemRepositoryInterface interface {
	GetAll() ([]models.Item, error)
	GetById(int64) (models.Item, error)
	GetByName(string) (models.Item, error)
	Create(*models.Item) (models.Item, error)
	Delete(int64) (int64, error)
	Update(int64, *models.Item) (models.Item, error)
	GetItemCategories(int64) ([]models.Category, error)
}

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{
		db: db,
	}
}

func (r *ItemRepository) GetAll() ([]models.Item, error) {
	items := make([]models.Item, 0)

	sqlStatement := `SELECT * FROM items`

	rows, queryErr := r.db.Query(sqlStatement)

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

func (r *ItemRepository) GetById(id int64) (models.Item, error) {
	var item models.Item

	sqlStatement := `SELECT * FROM items WHERE item_id = $1`

	row := r.db.QueryRow(sqlStatement, id)

	err := row.Scan(&item.ItemID, &item.Item, &item.Price)

	return item, err
}

func (r *ItemRepository) GetByName(name string) (models.Item, error) {
	var item models.Item

	sqlStatement := `SELECT * FROM items WHERE item = $1`

	row := r.db.QueryRow(sqlStatement, name)

	err := row.Scan(&item.ItemID, &item.Item, &item.Price)

	return item, err
}

func (r *ItemRepository) Create(itemReq *models.Item) (models.Item, error) {
	sqlStatement := `INSERT INTO items (item, price) VALUES ($1, $2) RETURNING *`

	var itemResp models.Item

	err := r.db.QueryRow(sqlStatement, itemReq.Item, itemReq.Price).Scan(&itemResp.ItemID, &itemResp.Item, &itemResp.Price)

	return itemResp, err
}

func (r *ItemRepository) Delete(id int64) (int64, error) {
	sqlStatement := `DELETE FROM items WHERE item_id = $1`

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

func (r *ItemRepository) Update(id int64, itemReq *models.Item) (models.Item, error) {
	sqlStatement := `UPDATE items SET item = $2, price = $3 WHERE item_id = $1 RETURNING *`

	var itemResp models.Item

	err := r.db.QueryRow(sqlStatement, id, itemReq.Item, itemReq.Price).Scan(&itemResp.ItemID, &itemResp.Item, &itemResp.Price)

	return itemResp, err
}

func (r *ItemRepository) GetItemCategories(id int64) ([]models.Category, error) {
	categories := make([]models.Category, 0)

	_, getErr := r.GetById(id)

	if getErr != nil {
		return nil, getErr
	}

	sqlStatement := `SELECT category_id, category FROM categories
	INNER JOIN categories_items
	USING (category_id)
	WHERE item_id = $1`

	rows, queryErr := r.db.Query(sqlStatement, id)

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
