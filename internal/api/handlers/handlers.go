package handlers

import "training/proj/internal/db/repositories"

type Handlers struct {
	CategoryHandler *CategoryHandler
	ItemHandler     *ItemHandler
	UserHandler     *UserHandler
}

func NewHandlers(cr *repositories.CategoryRepository, cir *repositories.CategoryItemRepository,
	ir *repositories.ItemRepository, ur *repositories.UserRepository) *Handlers {
	return &Handlers{
		CategoryHandler: NewCategoryHandler(cr, cir),
		ItemHandler:     NewItemHandler(ir),
		UserHandler:     NewUserHandler(ur),
	}
}
