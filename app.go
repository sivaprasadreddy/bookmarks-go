package main

import (
	"database/sql"
	"embed"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/sivaprasadreddy/bookmarks-go/bookmarks"
	"github.com/sivaprasadreddy/bookmarks-go/config"
	"github.com/sivaprasadreddy/bookmarks-go/database"
)

type App struct {
	Router             *mux.Router
	db                 *sql.DB
	bookmarkController *bookmarks.BookmarkController
}

func NewApp(config config.AppConfig) *App {
	app := &App{}
	app.init(config)
	return app
}

func (app *App) init(config config.AppConfig) {
	//logFile := initLogging()
	//defer logFile.Close()
	app.initLogging()

	app.db = database.GetDb(config)

	bookmarksRepo := bookmarks.NewBookmarkRepo(app.db)
	bookmarkSvc := bookmarks.NewBookmarkService(bookmarksRepo)
	app.bookmarkController = bookmarks.NewBookmarkController(bookmarkSvc)

	app.Router = app.setupRoutes()
}

//go:embed templates/*
var assetData embed.FS

//go:embed static
var staticFiles embed.FS

func (app *App) setupRoutes() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	var staticFS = http.FS(staticFiles)
	fs := http.FileServer(staticFS)

	router.PathPrefix("/static/").Handler(fs)

	router.HandleFunc("/", app.rootRouteHandler)
	apiRouter := router.PathPrefix("/api").Subrouter()
	app.setupBookmarkApiRoutes(apiRouter, app.bookmarkController)
	return router
}

func (app *App) rootRouteHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(assetData, "templates/index.html")
	if err != nil {
		log.Fatalf("error loading static assets: %v", err)
	}
	tmpl.Execute(w, nil)
}

func (app *App) setupBookmarkApiRoutes(router *mux.Router, controller *bookmarks.BookmarkController) {
	r := router.PathPrefix("/bookmarks").Subrouter()
	r.HandleFunc("", controller.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/{id:[0-9]+}", controller.GetById).Methods(http.MethodGet)
	r.HandleFunc("", controller.Create).Methods(http.MethodPost)
	r.HandleFunc("/{id:[0-9]+}", controller.Update).Methods(http.MethodPut)
	r.HandleFunc("/{id:[0-9]+}", controller.Delete).Methods(http.MethodDelete)
}

func (app *App) initLogging() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.InfoLevel)
}

func (app *App) initFileLogging() *os.File {
	logFile, err := os.OpenFile("bookmarks.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	log.SetOutput(logFile)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	return logFile
}
