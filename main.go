package main

import (
	"effective-mobile/app"
	"effective-mobile/logger"
	"effective-mobile/routes"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var port = os.Getenv("APP_PORT")

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

	logger.LogInfo(fmt.Sprintf("Starting application on port %s", port))

	// create router
	router := routes.InitRoutes(application)

	server := http.Server{
		Addr:         fmt.Sprintf("localhost:%s", port),
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
