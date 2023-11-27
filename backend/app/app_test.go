package app_test

import (
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

func (suite *AppTestSuite) SetupSuite() {
	suite.db.Conn.MustExec("BEGIN")
}

func (suite *AppTestSuite) TearDownSuite() {
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

func TestAppTesSuite(t *testing.T) {
	db := app.ConnectDb()
	suite.Run(t, &AppTestSuite{db: db, frontendPath: "./mocks_index.html"})
}
