package bookmarks

import (
	"github.com/gin-gonic/gin"
	"github.com/sivaprasadreddy/bookmarks-go/internal/config"
	"github.com/sivaprasadreddy/bookmarks-go/testsupport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ControllerTestSuite struct {
	suite.Suite
	PgContainer *testsupport.PostgresContainer
	cfg         config.AppConfig
	app         *App
	router      *gin.Engine
}

func (suite *ControllerTestSuite) SetupSuite() {
	suite.PgContainer = testsupport.InitPostgresContainer()
	suite.cfg = config.GetConfig(".env")
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
