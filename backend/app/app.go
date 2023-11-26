package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Application struct {
	Router *mux.Router
	// TODO: Add storage for messages
	messages []Message
}

type Message struct {
	Id        int    `json:"id"`
	ParentId  *int   `json:"parentId,omitempty"`
	Author    string `json:"author"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

type ParentId struct {
	Value int
	Valid bool
}

func NewApplication(r *mux.Router, frontendPath string) *Application {
	one := 1
	three := 3

	messages := []Message{
		{Id: 1, ParentId: nil, Author: "John", Title: "First!", Content: "Hello World!", CreatedAt: "2020-01-01 12:00:00"},
		{Id: 2, ParentId: nil, Author: "Jane", Title: "Second!", Content: "Hello John!", CreatedAt: "2020-01-01 12:00:01"},
		{Id: 3, ParentId: &one, Author: "Jane", Title: "", Content: "Hello John! Accidentally made a completely new post :)", CreatedAt: "2020-01-01 12:00:02"},
		{Id: 4, ParentId: &three, Author: "John", Title: "", Content: "I get it. The UX of this board is kind of terrible.", CreatedAt: "2020-01-01 12:00:03"},
		{Id: 5, ParentId: &one, Author: "Paul", Title: "", Content: "Hi there, John!", CreatedAt: "2020-01-01 12:00:04"},
	}

	myApp := &Application{Router: r, messages: messages}
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
	json.NewEncoder(w).Encode(a.messages)
}

func (a *Application) getHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
