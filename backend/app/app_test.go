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

const (
	NUM_MESSAGES_IN_DEFAULT_MIGRATION = 5
)

type AppTestSuite struct {
	suite.Suite
	db app.Database
	frontendPath string
}

func (suite *AppTestSuite) SetupSuite() {
	suite.db.Conn.MustExec("BEGIN");
}

func (suite *AppTestSuite) TearDownSuite() {
	suite.db.Conn.MustExec("ROLLBACK");
}

func (suite *AppTestSuite) Test_GetComments() {
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

func TestAppTesSuite(t *testing.T) {
	db := app.ConnectDb()
	suite.Run(t, &AppTestSuite{db: db, frontendPath: "./mocks_index.html"})
}
