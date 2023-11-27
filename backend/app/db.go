package app

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Note: Could be made an interface to allow for mocking in tests
// if one likes that kind of thing.
type Database struct {
	Conn *sqlx.DB
}

const (
	DB_USER_ENV     = "DB_USER"
	DB_PASSWORD_ENV = "DB_PASSWORD"
	DB_HOST_ENV     = "DB_HOST"
	DB_PORT_ENV     = "DB_PORT"
	DB_NAME_ENV     = "DB_NAME"
)

func ConnectDb() Database {
	dbUser := GetEnv(DB_USER_ENV)
	dbPassword := GetEnv(DB_PASSWORD_ENV)
	dbHost := GetEnv(DB_HOST_ENV)
	dbPort := GetEnv(DB_PORT_ENV)
	dbName := GetEnv(DB_NAME_ENV)

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	log.Printf("Connecting to database with connection %s", connectionString)
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")
	return Database{Conn: db}
}

func (db *Database) getMessages() ([]Message, error) {
	var messages []Message
	err := db.Conn.Select(&messages,
		`SELECT m.id, m.parent_id, a.username, m.title, m.content, m.created_at
		FROM message m
		JOIN author a ON m.author_id = a.id`)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s environment variable not set", key)
	}
	return value
}
