package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/sivaprasadreddy/bookmarks-go/config"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const postgresImage = "postgres:15.3-alpine"
const postgresPort = "5432"
const postgresUserName = "postgres"
const postgresPassword = "postgres"
const postgresDbName = "postgres"

type PostgresContainer struct {
	Container testcontainers.Container
	CloseFn   func()
	Host      string
	Port      string
	Database  string
	Username  string
	Password  string
}

// SetupPostgres creates an instance of the postgres container type
func SetupPostgres(ctx context.Context) (*PostgresContainer, error) {
	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage(postgresImage),
		postgres.WithDatabase(postgresDbName),
		postgres.WithUsername(postgresUserName),
		postgres.WithPassword(postgresPassword),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	host, _ := container.Host(ctx)
	hostPort, _ := container.MappedPort(ctx, postgresPort)

	return &PostgresContainer{
		Container: container,
		CloseFn: func() {
			if err := container.Terminate(ctx); err != nil {
				log.Fatalf("error terminating postgres container: %s", err)
			}
		},
		Host:     host,
		Port:     hostPort.Port(),
		Database: postgresDbName,
		Username: postgresUserName,
		Password: postgresPassword,
	}, nil
}

func overrideEnv(pgC *PostgresContainer) {
	os.Setenv("APP_DB_HOST", pgC.Host)
	os.Setenv("APP_DB_PORT", fmt.Sprint(pgC.Port))
	os.Setenv("APP_DB_USERNAME", pgC.Username)
	os.Setenv("APP_DB_PASSWORD", pgC.Password)
	os.Setenv("APP_DB_NAME", pgC.Database)
	os.Setenv("APP_DB_RUN_MIGRATIONS", "true")
}

var cfg config.AppConfig
var app *App
var router *gin.Engine

func TestMain(m *testing.M) {
	//Common Setup
	ctx := context.Background()
	pgContainer, err := SetupPostgres(ctx)
	if err != nil {
		log.Error("failed to setup Postgres container")
		panic(err)
	}
	defer pgContainer.CloseFn()
	overrideEnv(pgContainer)

	cfg = config.GetConfig()
	app = NewApp(cfg)
	router = app.Router

	code := m.Run()

	//Common Teardown
	os.Exit(code)
}

func TestGetAllBookmarks(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/bookmarks", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	actualResponseJson := response.Body.String()
	assert.NotEqual(t, "[]", actualResponseJson)

	if actualResponseJson == "[]" {
		t.Errorf("Expected an non-empty array. Got %s", actualResponseJson)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
