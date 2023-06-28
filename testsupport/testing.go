package testsupport

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
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
	panicIfError(os.Setenv("DB_HOST", pgC.Host))
	panicIfError(os.Setenv("DB_PORT", fmt.Sprint(pgC.Port)))
	panicIfError(os.Setenv("DB_USERNAME", pgC.Username))
	panicIfError(os.Setenv("DB_PASSWORD", pgC.Password))
	panicIfError(os.Setenv("DB_NAME", pgC.Database))
	panicIfError(os.Setenv("DB_RUN_MIGRATIONS", "true"))
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
