package main

import (
	"fmt"
	"log"
	"message-board-backend/app"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	FRONTEND_PATH_ENV = "MESSAGE_BOARD_FRONTEND_PATH"
	PORT_ENV          = "PORT"
	DB_USER_ENV       = "DB_USER"
	DB_PASSWORD_ENV   = "DB_PASSWORD"
	DB_HOST_ENV       = "DB_HOST"
	DB_PORT_ENV       = "DB_PORT"
	DB_NAME_ENV       = "DB_NAME"
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

	db := connectDb()
	application := app.NewApplication(r, frontendPath, db)
	log.Fatal(http.ListenAndServe(":"+port, application.Router))
}

func connectDb() *sqlx.DB {
	dbUser := getEnv(DB_USER_ENV)
	dbPassword := getEnv(DB_PASSWORD_ENV)
	dbHost := getEnv(DB_HOST_ENV)
	dbPort := getEnv(DB_PORT_ENV)
	dbName := getEnv(DB_NAME_ENV)

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	log.Printf("Connecting to database with connection %s", connectionString)
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")
	return db
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s environment variable not set", key)
	}
	return value
}
