package testsupport

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func InitPostgresContainer() *PostgresContainer {
	ctx := context.Background()
	pgContainer, err := SetupPostgres(ctx)
	if err != nil {
		log.Fatalf("failed to setup Postgres container")
		return nil
	}
	overrideEnv(pgContainer)
	return pgContainer
}

func overrideEnv(pgC *PostgresContainer) {
	os.Setenv("APP_DB_HOST", pgC.Host)
	os.Setenv("APP_DB_PORT", fmt.Sprint(pgC.Port))
	os.Setenv("APP_DB_USERNAME", pgC.Username)
	os.Setenv("APP_DB_PASSWORD", pgC.Password)
	os.Setenv("APP_DB_NAME", pgC.Database)
	os.Setenv("APP_DB_RUN_MIGRATIONS", "true")
}
