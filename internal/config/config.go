package config

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"training/proj/internal/api/handlers"
	"training/proj/internal/db/repositories"
)

type Config struct {
	Address   string
	JwtSecret string
	DSN       string
	DbName    string
}

func NewConfig() *Config {
	return &Config{}
}

func (cfg *Config) ParseFlags() error {
	address := fmt.Sprintf("%v:%v", os.Getenv("APP_HOST"), os.Getenv("APP_INTERNAL_PORT"))
	flag.StringVar(&cfg.Address, "address", address, "API server address")
	flag.StringVar(&cfg.JwtSecret, "jwtSecret", os.Getenv("JWT_SECRET_KEY"), "JWT secret key")

	connectionString := fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=%v",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("DB_INTERNAL_PORT"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSL_MODE"))
	flag.StringVar(&cfg.DSN, "DSN", connectionString, "DSN")
	flag.StringVar(&cfg.DbName, "dbName", os.Getenv("POSTGRES_DB"), "DB name")
	return nil
}

func (c *Config) InitializeHandlers(r *repositories.Repositories) *handlers.Handlers {
	return handlers.NewHandlers(r.CategoryRepository, r.CategoryItemRepository, r.ItemRepository, r.UserRepository)
}

func (c *Config) InitializeRepositories(db *sql.DB) *repositories.Repositories {
	return repositories.NewRepositories(db)
}
