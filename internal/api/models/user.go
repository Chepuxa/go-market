package models

type User struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email" validate:"required"`
	Username  string `json:"username" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Password  string `json:"password,omitempty" validate:"required"`
}
