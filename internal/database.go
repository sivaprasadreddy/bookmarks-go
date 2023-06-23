package bookmarks

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

func GetDb(config AppConfig) *pgx.Conn {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DbHost, config.DbPort, config.DbUserName, config.DbPassword, config.DbDatabase)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	if config.DbRunMigrations {
		runMigrations(config)
	}
	return conn
}

func runMigrations(config AppConfig) {
	sourceURL := config.DbMigrationsLocation
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.DbUserName, config.DbPassword, config.DbHost, config.DbPort, config.DbDatabase)
	log.Printf("DB Migration sourceURL: %s\n", sourceURL)
	log.Printf("DB Migration URL: %s\n", databaseURL)
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		log.Fatalf("Database migration error: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Databse migrate.up() error: %v", err)
	}
	log.Printf("Database migration completed")
}
