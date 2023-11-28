package app_test

import (
	"bytes"
	"encoding/json"
	"message-board-backend/app"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

// badness: magic numbers, but I do not have time to do proper test setup and
// teardown so I'm using the default migration for now.
const (
	NUM_ROOT_MESSAGES_IN_DEFAULT_MIGRATION = 2
	NUM_COMMENTS_ON_MESSAGE_1              = 2
)

type AppTestSuite struct {
	suite.Suite
	db           app.Database
	frontendPath string
}

func (suite *AppTestSuite) BeforeTest(suiteName, testName string) {
	suite.db.Conn.MustExec("BEGIN")
}

func (suite *AppTestSuite) AfterTest(suiteName, testName string) {
	suite.db.Conn.MustExec("ROLLBACK")
}

func (suite *AppTestSuite) Test_GetMessages() {
	myApp := app.NewApplication(mux.NewRouter(), suite.frontendPath, suite.db)
	method := "GET"
	path := "/api/messages"
	expectedStatus := http.StatusOK

	request, err := http.NewRequest(method, path, nil)
	suite.NoError(err)
	response := httptest.NewRecorder()
	myApp.Router.ServeHTTP(response, request)

	messages := []app.Message{}
	err = json.NewDecoder(response.Body).Decode(&messages)
	suite.NoError(err)

	suite.Len(messages, NUM_ROOT_MESSAGES_IN_DEFAULT_MIGRATION)
	suite.Equal(expectedStatus, response.Code)
}

func (suite *AppTestSuite) Test_GetComments() {
	myApp := app.NewApplication(mux.NewRouter(), suite.frontendPath, suite.db)
	method := "GET"
	path := "/api/messages/1/comments"
	expectedStatus := http.StatusOK

	request, err := http.NewRequest(method, path, nil)
	suite.NoError(err)
	response := httptest.NewRecorder()
	myApp.Router.ServeHTTP(response, request)

	messages := []app.Message{}
	err = json.NewDecoder(response.Body).Decode(&messages)
	suite.NoError(err)

	suite.Len(messages, NUM_COMMENTS_ON_MESSAGE_1)
	suite.Equal(expectedStatus, response.Code)
}

func (suite *AppTestSuite) Test_GetComments_InvalidParentId() {
	myApp := app.NewApplication(mux.NewRouter(), suite.frontendPath, suite.db)
	method := "GET"
	path := "/api/messages/invalid/comments"
	expectedStatus := http.StatusBadRequest

	request, err := http.NewRequest(method, path, nil)
	suite.NoError(err)

	response := httptest.NewRecorder()
	myApp.Router.ServeHTTP(response, request)

	suite.Equal(expectedStatus, response.Code)
}

func (suite *AppTestSuite) Test_GetComments_NoSuchMessage() {
	myApp := app.NewApplication(mux.NewRouter(), suite.frontendPath, suite.db)
	method := "GET"
	path := "/api/messages/-1/comments"

	// FIXME: this should be http.StatusNotFound, but due to the current
	// implementation it isn't
	expectedStatus := http.StatusOK

	request, err := http.NewRequest(method, path, nil)
	suite.NoError(err)
	response := httptest.NewRecorder()
	myApp.Router.ServeHTTP(response, request)

	suite.Equal(expectedStatus, response.Code)
}

func (suite *AppTestSuite) Test_CreateMessage() {
	myApp := app.NewApplication(mux.NewRouter(), suite.frontendPath, suite.db)
	method := "POST"
	path := "/api/messages"
	expectedStatus := http.StatusCreated

	message := app.InputMessage{
		Author:  "John",
		Title:   "Hello",
		Content: "World",
	}

	requestBody, err := json.Marshal(message)
	suite.NoError(err)

	request, err := http.NewRequest(method, path, bytes.NewBuffer(requestBody))
	suite.NoError(err)
	response := httptest.NewRecorder()
	myApp.Router.ServeHTTP(response, request)

	var returnedMessage app.Message
	err = json.NewDecoder(response.Body).Decode(&returnedMessage)
	suite.NoError(err)

	suite.Equal(expectedStatus, response.Code)
	suite.Greater(returnedMessage.Id, NUM_ROOT_MESSAGES_IN_DEFAULT_MIGRATION)
	suite.Equal(message.Title, returnedMessage.Title)
	suite.Equal(message.Content, returnedMessage.Content)
	suite.Equal(message.Author, returnedMessage.Author)
	suite.Nil(returnedMessage.ParentId)
}

func (suite *AppTestSuite) Test_CreateComment() {
	myApp := app.NewApplication(mux.NewRouter(), suite.frontendPath, suite.db)
	method := "POST"
	path := "/api/messages"
	expectedStatus := http.StatusCreated
	one := 1

	message := app.InputMessage{
		Author:   "John",
		Content:  "World",
		ParentId: &one,
	}

	requestBody, err := json.Marshal(message)
	suite.NoError(err)

	request, err := http.NewRequest(method, path, bytes.NewBuffer(requestBody))
	suite.NoError(err)
	response := httptest.NewRecorder()
	myApp.Router.ServeHTTP(response, request)

	var returnedMessage app.Message
	err = json.NewDecoder(response.Body).Decode(&returnedMessage)
	suite.NoError(err)

	suite.Equal(expectedStatus, response.Code)
	suite.Greater(returnedMessage.Id, NUM_ROOT_MESSAGES_IN_DEFAULT_MIGRATION)
	suite.Equal(message.ParentId, returnedMessage.ParentId)
}

func (suite *AppTestSuite) Test_DeleteMessage() {
	myApp := app.NewApplication(mux.NewRouter(), suite.frontendPath, suite.db)
	method := "DELETE"
	path := "/api/messages/1"
	expectedStatus := http.StatusOK

	request, err := http.NewRequest(method, path, nil)
	suite.NoError(err)
	response := httptest.NewRecorder()
	myApp.Router.ServeHTTP(response, request)

	var message app.Message
	err = json.NewDecoder(response.Body).Decode(&message)
	suite.NoError(err)

	suite.Equal(response.Code, expectedStatus)
	suite.Equal(app.REDACTED_USERNAME, message.Author)
	suite.Equal(app.REDACTED_TITLE, message.Title)
	suite.Equal(app.REDACTED_CONTENT, message.Content)
}

func (suite *AppTestSuite) Test_DeleteMessage_NoMessage() {
	myApp := app.NewApplication(mux.NewRouter(), suite.frontendPath, suite.db)
	method := "DELETE"
	path := "/api/messages/99"
	expectedStatus := http.StatusNotFound

	request, err := http.NewRequest(method, path, nil)
	suite.NoError(err)
	response := httptest.NewRecorder()
	myApp.Router.ServeHTTP(response, request)

	suite.Equal(expectedStatus, response.Code)
}

func TestAppTesSuite(t *testing.T) {
	db := app.ConnectDb()
	suite.Run(t, &AppTestSuite{db: db, frontendPath: "./mocks_index.html"})
}
