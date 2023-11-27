package main

import (
	"github.com/gorilla/mux"
	"log"
	"message-board-backend/app"
	"net/http"
	"os"
)

const (
	FRONTEND_PATH_ENV = "MESSAGE_BOARD_FRONTEND_PATH"
)

func main() {
	r := mux.NewRouter()

	frontendPath := os.Getenv(FRONTEND_PATH_ENV)
	if frontendPath == "" {
		log.Fatalf("%s environment variable not set", FRONTEND_PATH_ENV)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT environment variable not set")
	}

	log.Printf("Serving frontend from %s on port %s", frontendPath, port)

	application := app.NewApplication(r, frontendPath)
	log.Fatal(http.ListenAndServe(":" + port, application.Router))
}
