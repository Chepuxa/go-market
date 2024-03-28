package models

type Category struct {
	CategoryID int64  `json:"category_id"`
	Category   string `json:"category" validate:"required"`
}
