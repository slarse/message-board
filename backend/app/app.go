package app

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"log"

	"github.com/gorilla/mux"
)

type Application struct {
	Router *mux.Router
	db     Database
}

type InputMessage struct {
	ParentId *int   `json:"parentId,omitempty" db:"parent_id"`
	Author   string `json:"author" db:"username"`
	Title    string `json:"title" db:"title"`
	Content  string `json:"content" db:"content"`
}

type Message struct {
	InputMessage
	Id        int    `json:"id" db:"id"`
	CreatedAt string `json:"createdAt" db:"created_at"`
}

func NewApplication(r *mux.Router, frontendPath string, db Database) *Application {
	myApp := &Application{Router: r, db: db}
	myApp.setupRoutes(frontendPath)
	return myApp
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s environment variable not set", key)
	}
	return value
}

func (a *Application) setupRoutes(frontendPath string) {
	a.Router.HandleFunc("/api/health", a.getHealth).Methods("GET")
	a.Router.HandleFunc("/api/messages", a.getMessages).Methods("GET")
	a.Router.HandleFunc("/api/messages", a.createMessage).Methods("POST")
	a.Router.HandleFunc("/api/messages/{messageId}", a.deleteMessage).Methods("DELETE")
	a.Router.HandleFunc("/api/messages/{messageId}/comments", a.getComments)

	// Normally, we would serve the frontend from a static file server (e.g. an S3 bucket)
	// with a CDN in front of it (e.g. CloudFront). But for this example we will
	// serve it from the backend itself for the sake of less moving parts.
	frontendFileServer := http.FileServer(http.Dir(frontendPath))
	a.Router.PathPrefix("/").Handler(frontendFileServer)
}

func (a *Application) getMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := a.db.getMessages(nil)
	if err != nil {
		log.Printf("Error getting messages from database: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messages)
}

func (a *Application) getComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageId, err := strconv.ParseInt(vars["messageId"], 10, 64)
	if err != nil {
		log.Printf("Invalid messageId: %s", err)
		http.Error(w, "Invalid messageId", http.StatusBadRequest)
		return
	}

	comments, err := a.db.getMessages(&messageId)
	if err != nil {
		log.Printf("Error getting comments from database: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comments)
}

func (a *Application) createMessage(w http.ResponseWriter, r *http.Request) {
	var inputMessage InputMessage
	err := json.NewDecoder(r.Body).Decode(&inputMessage)
	if err != nil {
		log.Printf("Error decoding input message: %s", err)
		http.Error(w, "Invalid input message", http.StatusBadRequest)
		return
	}

	created, err := a.db.createMessage(inputMessage)
	if err != nil {
		log.Printf("Error creating message: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (a *Application) deleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageId, err := strconv.ParseInt(vars["messageId"], 10, 64)
	if err != nil {
		log.Printf("Invalid messageId: %s", err)
		http.Error(w, "Invalid messageId", http.StatusBadRequest)
		return
	}

	message, err := a.db.deleteMessage(messageId)
	if err != nil && err.Error() == "sql: no rows in result set" {
		log.Printf("No message for id: %d", messageId)
		http.NotFound(w, r)
		return
	} else if err != nil {
		log.Printf("Error deleting message: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}

func (a *Application) getHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
