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
	NUM_MESSAGES_IN_DEFAULT_MIGRATION = 5
	NUM_MESSAGES_ROOTED_IN_MESSAGE_1  = 4
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

	suite.Len(messages, NUM_MESSAGES_IN_DEFAULT_MIGRATION)
	suite.Equal(expectedStatus, response.Code)
}

func (suite *AppTestSuite) Test_GetMessagesByRootMessageId() {
	myApp := app.NewApplication(mux.NewRouter(), suite.frontendPath, suite.db)
	method := "GET"
	path := "/api/messages/1"
	expectedStatus := http.StatusOK

	request, err := http.NewRequest(method, path, nil)
	suite.NoError(err)
	response := httptest.NewRecorder()
	myApp.Router.ServeHTTP(response, request)

	messages := []app.Message{}
	err = json.NewDecoder(response.Body).Decode(&messages)
	suite.NoError(err)

	suite.Len(messages, NUM_MESSAGES_ROOTED_IN_MESSAGE_1)
	suite.Equal(expectedStatus, response.Code)
}

func (suite *AppTestSuite) Test_GetMessagesByRootMessageId_InvalidId() {
	myApp := app.NewApplication(mux.NewRouter(), suite.frontendPath, suite.db)
	method := "GET"
	path := "/api/messages/invalid"
	expectedStatus := http.StatusBadRequest

	request, err := http.NewRequest(method, path, nil)
	suite.NoError(err)

	response := httptest.NewRecorder()
	myApp.Router.ServeHTTP(response, request)

	suite.Equal(expectedStatus, response.Code)
}

func (suite *AppTestSuite) Test_GetMessagesByRootMessageId_NoMessages() {
	myApp := app.NewApplication(mux.NewRouter(), suite.frontendPath, suite.db)
	method := "GET"
	path := "/api/messages/-1"
	expectedStatus := http.StatusNotFound

	request, err := http.NewRequest(method, path, nil)
	suite.NoError(err)
	response := httptest.NewRecorder()
	myApp.Router.ServeHTTP(response, request)

	suite.Equal(expectedStatus, response.Code)
}

func (suite *AppTestSuite) Test_CreateRootMessage() {
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
	suite.Greater(returnedMessage.Id, NUM_MESSAGES_IN_DEFAULT_MIGRATION)
	suite.Equal(returnedMessage.Title, message.Title)
	suite.Equal(returnedMessage.Content, message.Content)
	suite.Equal(returnedMessage.Author, message.Author)
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
	suite.Greater(returnedMessage.Id, NUM_MESSAGES_IN_DEFAULT_MIGRATION)
	suite.Equal(returnedMessage.ParentId, message.ParentId)
}

func TestAppTesSuite(t *testing.T) {
	db := app.ConnectDb()
	suite.Run(t, &AppTestSuite{db: db, frontendPath: "./mocks_index.html"})
}
