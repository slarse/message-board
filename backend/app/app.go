package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Application struct {
	Router *mux.Router
	// TODO: Add storage for comments
}

func NewApplication(r *mux.Router) *Application {
	myApp := &Application{Router: r}
	myApp.setupRoutes()
	return myApp
}

func (a *Application) setupRoutes() {
	// TODO: Setup routes and link to functions
}

// Template function, needs to be linked to a route
func (a *Application) getComments(w http.ResponseWriter, r *http.Request) {
	// TODO: Get comments from database...
}
