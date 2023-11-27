package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
)

type Application struct {
	Router   *mux.Router
	Database *sqlx.DB
}

type Message struct {
	Id        int    `json:"id" db:"id"`
	ParentId  *int   `json:"parentId,omitempty" db:"parent_id"`
	Author    string `json:"author" db:"username"`
	Title     string `json:"title" db:"title"`
	Content   string `json:"content" db:"content"`
	CreatedAt string `json:"createdAt" db:"created_at"`
}

type ParentId struct {
	Value int
	Valid bool
}

func NewApplication(r *mux.Router, frontendPath string, db *sqlx.DB) *Application {
	myApp := &Application{Router: r, Database: db}
	myApp.setupRoutes(frontendPath)
	return myApp
}

func (a *Application) setupRoutes(frontendPath string) {
	a.Router.HandleFunc("/api/health", a.getHealth).Methods("GET")
	a.Router.HandleFunc("/api/messages", a.getMessages).Methods("GET")

	// Normally, we would serve the frontend from a static file server (e.g. an S3 bucket)
	// with a CDN in front of it (e.g. CloudFront). But for this example we will
	// serve it from the backend itself for the sake of less moving parts.
	frontendFileServer := http.FileServer(http.Dir(frontendPath))
	a.Router.PathPrefix("/").Handler(frontendFileServer)
}

func (a *Application) getMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := a.getMessagesFromDb()
	if err != nil {
		log.Printf("Error getting messages from database: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messages)
}

func (a *Application) getHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func (a *Application) getMessagesFromDb() ([]Message, error) {
	var messages []Message
	err := a.Database.Select(&messages,
		`SELECT m.id, m.parent_id, a.username, m.title, m.content, m.created_at
		FROM message m
		JOIN author a ON m.author_id = a.id`)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
