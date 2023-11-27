package main

import (
	"log"
	"message-board-backend/app"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	FRONTEND_PATH_ENV = "MESSAGE_BOARD_FRONTEND_PATH"
	PORT_ENV          = "PORT"
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

	db := app.ConnectDb()
	application := app.NewApplication(r, frontendPath, db)
	log.Fatal(http.ListenAndServe(":"+port, application.Router))
}
