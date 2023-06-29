package db

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/sivaprasadreddy/bookmarks-go/internal/config"
	"github.com/sivaprasadreddy/bookmarks-go/internal/logging"
)

func GetDb(config config.AppConfig, logger *logging.Logger) *pgx.Conn {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DbHost, config.DbPort, config.DbUserName, config.DbPassword, config.DbDatabase)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		logger.Fatal(err)
	}
	if config.DbRunMigrations {
		runMigrations(config, logger)
	}
	return conn
}

func runMigrations(config config.AppConfig, logger *logging.Logger) {
	sourceURL := config.DbMigrationsLocation
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.DbUserName, config.DbPassword, config.DbHost, config.DbPort, config.DbDatabase)
	logger.Infof("DB Migration sourceURL: %s\n", sourceURL)
	logger.Infof("DB Migration URL: %s\n", databaseURL)
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		logger.Fatalf("Database migration error: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatalf("Databse migrate.up() error: %v", err)
	}
	logger.Infof("Database migration completed")
}
