package main

import (
	"fmt"
	bookmarks "github.com/sivaprasadreddy/bookmarks-go/internal"
	"github.com/sivaprasadreddy/bookmarks-go/internal/config"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.GetConfig(".env")
	app := bookmarks.NewApp(cfg)

	port := fmt.Sprintf(":%d", cfg.AppPort)
	srv := &http.Server{
		Handler:        app.Router,
		Addr:           port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("listening on port %d", cfg.AppPort)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
