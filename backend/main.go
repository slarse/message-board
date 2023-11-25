package main

import (
	"backend-test-template/app"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	application := app.NewApplication(r)
	log.Fatal(http.ListenAndServe(":8080", application.Router))
}
