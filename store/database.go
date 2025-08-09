package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		dbHost,
		dbUser,
		dbPassword,
		dbName,
	)

	db, err := sql.Open("pgx", connectionString)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *sql.DB, dir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose: set dialect %v", err)
	}

	if err := goose.Up(db, dir); err != nil {
		return fmt.Errorf("goose: migrate up %v", err)
	}

	return nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)

	defer func() {
		goose.SetBaseFS(nil)
	}()

	return Migrate(db, dir)
}

func SetupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")

	if err != nil {
		t.Fatalf("open db failed: %v", err)
	}

	// run the migrations
	if err := Migrate(db, "../migrations"); err != nil {
		t.Fatalf("migrating test db failed: %v", err)
	}

	if _, err := db.Exec(`BEGIN TRANSACTION; DELETE FROM users; DELETE FROM user_subscriptions; COMMIT`); err != nil {
		t.Fatalf("truncating db failed: %v", err)
	}

	return db
}
