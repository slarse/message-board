package app_test

import (
	"message-board-backend/app"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_GetComments(t *testing.T) {
	myApp := app.NewApplication(mux.NewRouter())
	method := // TODO: Specify the method
	path := // TODO: Specify the path
	expectedStatus := // TODO: Specify the expected status

	request, err := http.NewRequest(method, path, nil)
	assert.NoError(t, err)
	response := httptest.NewRecorder()
	myApp.Router.ServeHTTP(response, request)

	comments := []app.Comment{}
	err = json.NewDecoder(response.Body).Decode(&comments)
	assert.NoError(t, err)

	// Assert expected results (maybe add more assertions?)
	assert.Len(t, comments, 0)
	assert.Equal(t, expectedStatus, response.Code)
}

func Test_CreateComment(t *testing.T) {
	myApp := app.NewApplication(mux.NewRouter())
	method := // TODO: Specify the method
	path := // TODO: Specify the path
	expectedStatus := // TODO: Specify the expected status
	body := //TODO: Specify the body

	request, err := http.NewRequest(method, path, bytes.NewReader(body))
	assert.NoError(t, err)
	response := httptest.NewRecorder()
	myApp.Router.ServeHTTP(response, request)

	comments := []app.Comment{}
	err = json.NewDecoder(response.Body).Decode(&comments)
	assert.NoError(t, err)

	// Assert expected results (maybe add more assertions?)
	assert.Len(t, comments, 1)
	assert.Equal(t, expectedStatus, response.Code)
}

// TODO: Write test for getComment, updateComment, deleteComment, ect...
