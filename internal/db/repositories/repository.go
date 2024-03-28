package repositories

import "database/sql"

type Repositories struct {
	CategoryRepository     *CategoryRepository
	ItemRepository         *ItemRepository
	UserRepository         *UserRepository
	CategoryItemRepository *CategoryItemRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		CategoryRepository:     NewCategoryRepository(db),
		ItemRepository:         NewItemRepository(db),
		UserRepository:         NewUserRepository(db),
		CategoryItemRepository: NewCategoryItemRepository(db),
	}
}
