package main

import (
	"time"
	"training/proj/internal/api"
	"training/proj/internal/config"
	"training/proj/internal/db"
	"training/proj/internal/logger"
	"training/proj/internal/scheduler"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	defer logger.CloseLogger()

	err := godotenv.Load()
	if err != nil {
		logger.Logger.Fatal("Error loading .env file", zap.Error(err))
	}

	cfg := config.NewConfig()

	err = cfg.ParseFlags()
	if err != nil {
		logger.Logger.Fatal("Failed to parse command-line flags", zap.Error(err))
	}

	conn, err := db.Connect(cfg)
	if err != nil {
		logger.Logger.Fatal("Failed to connect to the database", zap.Error(err))
		panic(err)
	}
	defer conn.Close()

	migrationErr := db.CreateTables(conn, cfg)
	if migrationErr != nil {
		logger.Logger.Fatal("Failed to apply migrations", zap.Error(err))
		panic(migrationErr)
	}

	repositories := cfg.InitializeRepositories(conn)
	handlers := cfg.InitializeHandlers(repositories)
	srv := api.NewAPI(logger.Logger, cfg, handlers)

	sch := scheduler.NewScheduler(repositories, logger.Logger, srv.Wg)
	go func() {
		for {
			sch.ExternalDbFill()
			time.Sleep(1 * time.Hour)
		}
	}()

	err = srv.Run()
	if err != nil {
		logger.Logger.Fatal("Failed to start the server", zap.Error(err))
	}
}
