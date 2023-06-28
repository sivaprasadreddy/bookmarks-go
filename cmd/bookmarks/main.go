package main

import (
	"flag"
	"log"

	bookmarks "github.com/sivaprasadreddy/bookmarks-go/internal"
	"github.com/sivaprasadreddy/bookmarks-go/internal/config"
)

func main() {
	var confFile string
	flag.StringVar(&confFile, "conf", ".env", "config path, eg: -conf app.dev")
	flag.Parse()
	cfg, err := config.GetConfig(confFile)
	if err != nil {
		log.Fatal(err)
	}
	app := bookmarks.NewApp(cfg)
	app.Run()
}
