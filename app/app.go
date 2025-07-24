package app

import (
	"database/sql"
	"effective-mobile/api"
	"effective-mobile/logger"
	"effective-mobile/migrations"
	"effective-mobile/store"
)

type Application struct {
	DB                   *sql.DB
	UsersHandler         *api.UsersHandler
	SubscriptionsHandler *api.SubscriptionsHandler
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()

	if err != nil {
		logger.LogError("Failed to open postgres connection", err)
		return nil, err
	}

	if err = store.MigrateFS(pgDB, migrations.FS, "."); err != nil {
		logger.LogError("Failed to complete migrations", err)
		return nil, err
	}

	// stores
	usersStore := store.NewPostgresUsersStore(pgDB)
	subscriptionsStore := store.NewPostgresSubscriptionsStore(pgDB)

	// handlers
	usersHandler := api.NewUsersHandler(usersStore)
	subscriptionsHandler := api.NewSubscriptionsHandler(subscriptionsStore)

	app := &Application{
		UsersHandler:         usersHandler,
		SubscriptionsHandler: subscriptionsHandler,
		DB:                   pgDB,
	}

	return app, nil
}
