package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UsersStore interface {
	CreateUser(user *User) error
	GetUserByID(id uuid.UUID) (*User, error)
}

type PostgresUsersStore struct {
	db *sql.DB
}

func NewPostgresUsersStore(db *sql.DB) *PostgresUsersStore {
	return &PostgresUsersStore{
		db: db,
	}
}

func (store *PostgresUsersStore) CreateUser(user *User) error {
	query := `
		INSERT INTO users (first_name, last_name)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	if err := store.db.QueryRow(query, user.FirstName, user.LastName).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (store *PostgresUsersStore) GetUserByID(id uuid.UUID) (*User, error) {
	user := &User{}

	query := `
		SELECT id, first_name, last_name, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := store.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}
