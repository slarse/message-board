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

	application := app.NewApplication(r, frontendPath)
	log.Fatal(http.ListenAndServe(":8080", application.Router))
}
