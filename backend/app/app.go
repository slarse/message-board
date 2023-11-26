package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Application struct {
	Router *mux.Router
	// TODO: Add storage for messages
}

type Message struct {
	Id        int    `json:"id"`
	ParentId  int    `json:"parentId"`
	Author    string `json:"author"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

func NewApplication(r *mux.Router, frontendPath string) *Application {
	myApp := &Application{Router: r}
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
	messages := []Message{
		{Id: 2, ParentId: 1, Author: "Jane", Title: "Second!", Content: "Hello John!", CreatedAt: "2020-01-01 12:00:01"},
		{Id: 1, ParentId: 0, Author: "John", Title: "First!", Content: "Hello World!", CreatedAt: "2020-01-01 12:00:00"},
	}
	json.NewEncoder(w).Encode(messages)
}

func (a *Application) getHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
