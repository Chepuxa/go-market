package repositories

import (
	"database/sql"
	"errors"
	"training/proj/internal/api/models"

	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepositoryInterface interface {
	Create(*models.User, []byte) (models.User, *pgconn.PgError)
	GetByEmail(string) (models.User, error)
	GetByUsername(string) (models.User, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(userReq *models.User, hashedPassword []byte) (models.User, *pgconn.PgError) {
	var userResp models.User

	sqlStatement := `INSERT INTO users (email, first_name, last_name, password, username)
	VALUES ($1, $2, $3, $4, $5) RETURNING user_id, email, first_name, last_name, username`

	row := r.db.QueryRow(sqlStatement,
		userReq.Email,
		userReq.FirstName,
		userReq.LastName,
		hashedPassword,
		userReq.Username)
	err := row.Scan(
		&userResp.UserID,
		&userResp.Email,
		&userResp.FirstName,
		&userResp.LastName,
		&userResp.Username)

	var pgErr *pgconn.PgError
	errors.As(err, &pgErr)

	return userResp, pgErr
}

func (r *UserRepository) GetByEmail(email string) (models.User, error) {
	var userResp models.User

	sqlStatement := `SELECT * FROM users WHERE email = $1`

	row := r.db.QueryRow(sqlStatement, email)
	err := row.Scan(
		&userResp.UserID,
		&userResp.Email,
		&userResp.FirstName,
		&userResp.LastName,
		&userResp.Password,
		&userResp.Username)

	return userResp, err
}

func (r *UserRepository) GetByUsername(username string) (models.User, error) {
	var userResp models.User

	sqlStatement := `SELECT * FROM users WHERE username = $1`

	row := r.db.QueryRow(sqlStatement, username)
	err := row.Scan(
		&userResp.UserID,
		&userResp.Email,
		&userResp.FirstName,
		&userResp.LastName,
		&userResp.Password,
		&userResp.Username)

	return userResp, err
}
