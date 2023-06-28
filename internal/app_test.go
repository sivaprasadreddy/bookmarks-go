package bookmarks

import (
	"encoding/json"
	"github.com/sivaprasadreddy/bookmarks-go/internal/config"
	"github.com/sivaprasadreddy/bookmarks-go/internal/domain"
	"github.com/sivaprasadreddy/bookmarks-go/testsupport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type ControllerTestSuite struct {
	suite.Suite
	PgContainer *testsupport.PostgresContainer
	cfg         config.AppConfig
	app         *App
	router      http.Handler
}

func (suite *ControllerTestSuite) SetupSuite() {
	suite.PgContainer = testsupport.InitPostgresContainer()
	cfg, err := config.GetConfig(".env")
	if err != nil {
		log.Fatal(err)
	}
	suite.cfg = cfg

	suite.app = NewApp(suite.cfg)
	suite.router = suite.app.Router
}

func (suite *ControllerTestSuite) TearDownSuite() {
	suite.PgContainer.CloseFn()
}

func TestControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}

func (suite *ControllerTestSuite) TestGetAllBookmarks() {
	t := suite.T()
	req, _ := http.NewRequest("GET", "/api/bookmarks", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	actualResponseJson := w.Body.String()
	assert.NotEqual(t, "[]", actualResponseJson)
}

func (suite *ControllerTestSuite) TestGetBookmarkById() {
	t := suite.T()
	req, _ := http.NewRequest(http.MethodGet, "/api/bookmarks/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.Bookmark
	err := json.NewDecoder(w.Body).Decode(&response)

	assert.Nil(t, err)
	assert.NotNil(t, response.Id)
	assert.Equal(t, "How To Remove Docker Containers, Images, Volumes, and Networks", response.Title)
	assert.Equal(t, "https://linuxize.com/post/how-to-remove-docker-images-containers-volumes-and-networks/", response.Url)
	assert.NotNil(t, response.CreatedDate)
}

func (suite *ControllerTestSuite) TestCreateBookmark() {
	t := suite.T()
	reqBody := strings.NewReader(`
		{
			"title": "Test Post title",
			"url":     "https://example.com"
		}
	`)

	req, _ := http.NewRequest(http.MethodPost, "/api/bookmarks", reqBody)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response domain.Bookmark
	err := json.NewDecoder(w.Body).Decode(&response)

	assert.Nil(t, err)
	assert.NotNil(t, response.Id)
	assert.Equal(t, "Test Post title", response.Title)
	assert.Equal(t, "https://example.com", response.Url)
	assert.NotNil(t, response.CreatedDate)
	assert.Nil(t, response.UpdatedDate)
}

func (suite *ControllerTestSuite) TestUpdateBookmark() {
	t := suite.T()
	reqBody := strings.NewReader(`
		{
			"title": "Test Updated title",
			"url":   "https://example2.com"
		}
	`)

	req, _ := http.NewRequest(http.MethodPut, "/api/bookmarks/1", reqBody)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.Bookmark
	err := json.NewDecoder(w.Body).Decode(&response)

	assert.Nil(t, err)
	assert.Equal(t, 1, response.Id)
	assert.Equal(t, "Test Updated title", response.Title)
	assert.Equal(t, "https://example2.com", response.Url)
	assert.NotNil(t, response.CreatedDate)
	assert.NotNil(t, response.UpdatedDate)
}

func (suite *ControllerTestSuite) TestDeleteBookmark() {
	t := suite.T()

	req, _ := http.NewRequest(http.MethodDelete, "/api/bookmarks/2", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
