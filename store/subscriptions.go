package store

import (
	"database/sql"
	"effective-mobile/utils"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID             uuid.UUID         `json:"id"`
	Name           string            `json:"service_name"`
	Price          uint16            `json:"price"`
	UserID         uuid.UUID         `json:"user_id"`
	StartDate      utils.CustomDate  `json:"start_date"`
	ExpirationDate *utils.CustomDate `json:"expiration_date"`
	CreatedAt      *time.Time        `json:"created_at"`
	UpdatedAt      *time.Time        `json:"updated_at"`
}

type SubscriptionsStore interface {
	CreateUserSubscription(subscription *Subscription) error
	GetUserSubscriptions(userID uuid.UUID, params *RequestParams) ([]*Subscription, error)
	GetTotalSubscriptionCost(userID uuid.UUID, params *RequestParams) (uint, error)
}

type PostgresSubscriptionsStore struct {
	db *sql.DB
}

func NewPostgresSubscriptionsStore(db *sql.DB) *PostgresSubscriptionsStore {
	return &PostgresSubscriptionsStore{
		db: db,
	}
}

func (store *PostgresSubscriptionsStore) CreateUserSubscription(subscription *Subscription) error {
	query := `
		INSERT INTO user_subscriptions (name, price, start_date, expiration_date, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	if err := store.db.QueryRow(
		query,
		subscription.Name,
		subscription.Price,
		subscription.StartDate,
		subscription.ExpirationDate,
		subscription.UserID,
	).Scan(&subscription.ID, &subscription.CreatedAt, &subscription.UpdatedAt); err != nil {
		return err
	}

	return nil
}

type RequestParams struct {
	ServiceName string
	StartDate   utils.CustomDate
	EndDate     utils.CustomDate
}

func (store *PostgresSubscriptionsStore) GetUserSubscriptions(userID uuid.UUID, params *RequestParams) ([]*Subscription, error) {
	subscriptions := []*Subscription{}

	query := `
		SELECT
			id,
			user_id, 
			name,
			price,
		    start_date,
			created_at,
			updated_at
		FROM user_subscriptions
		WHERE user_id = $1 AND name = $2 AND (
			start_date >= $3 AND start_date < $4
		)
	`

	rows, err := store.db.Query(query, userID, params.ServiceName, params.StartDate, params.EndDate)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		subscription := &Subscription{}
		if err := rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&subscription.Name,
			&subscription.Price,
			&subscription.StartDate,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

func (store *PostgresSubscriptionsStore) GetTotalSubscriptionCost(userID uuid.UUID, params *RequestParams) (uint, error) {
	var totalCost uint

	query := `
		SELECT
			SUM(price)
		FROM user_subscriptions
		WHERE user_id = $1 AND name = $2 AND (
			start_date >= $3 AND start_date < $4
		)
	`

	err := store.db.QueryRow(query, userID, params.ServiceName, params.StartDate, params.EndDate).Scan(&totalCost)

	if err != nil {
		return 0, err
	}

	return totalCost, nil
}
