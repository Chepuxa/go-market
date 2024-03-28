package models

type Item struct {
	ItemID int64  `json:"item_id"`
	Item   string `json:"item" validate:"required"`
	Price  int64  `json:"price" validate:"required"`
}
