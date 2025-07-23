package main

import (
	"effective-mobile/app"
	"effective-mobile/logger"
	"effective-mobile/routes"
	"fmt"
	"net/http"
	"time"
)

const PORT = 3000

func main() {
	application, err := app.NewApplication()

	if err != nil {
		logger.LogError("App creation failed", err)

		panic(err)
	}

	defer func() {
		if err := application.DB.Close(); err != nil {
			logger.LogFatal(err)
		}
	}()

	logger.LogInfo(fmt.Sprintf("Starting application on port %d", PORT))

	// create router
	router := routes.InitRoutes(application)

	server := http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", PORT),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      router,
	}

	err = server.ListenAndServe()

	if err != nil {
		logger.LogFatal(err)
	}
}
