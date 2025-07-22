package routes

import (
	"effective-mobile/app"
	"net/http"
)

func InitRoutes(app *app.Application) *http.ServeMux {
	router := &http.ServeMux{}

	router.HandleFunc("POST /users", app.UsersHandler.HandleCreateUser)
	router.HandleFunc("GET /users/{id}", app.UsersHandler.HandleGetUserByID)
	router.HandleFunc("POST /subscriptions", app.SubscriptionsHandler.HandleCreateSubscription)
	router.HandleFunc("GET /subscriptions", app.SubscriptionsHandler.HandleGetUserSubscriptions)
	router.HandleFunc("GET /subscriptions/cost", app.SubscriptionsHandler.HandleGetTotalSubscriptionCost)

	return router
}
