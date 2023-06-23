package bookmarks

import (
	"github.com/sivaprasadreddy/bookmarks-go/assets"
	"html/template"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

type App struct {
	Router             *gin.Engine
	db                 *pgx.Conn
	bookmarkController *BookmarkController
}

func NewApp(config AppConfig) *App {
	app := &App{}
	app.init(config)
	return app
}

func (app *App) init(config AppConfig) {
	//logFile := initLogging()
	//defer logFile.Close()
	app.initLogging()

	app.db = GetDb(config)

	bookmarksRepo := NewBookmarkRepo(app.db)
	app.bookmarkController = NewBookmarkController(bookmarksRepo)

	app.Router = app.setupRoutes()
}

func (app *App) setupRoutes() *gin.Engine {
	r := gin.Default()

	r.Any("/", app.rootRouteHandler)
	r.GET("/static/*filepath", func(c *gin.Context) {
		c.FileFromFS(path.Join("/", c.Request.URL.Path), http.FS(assets.StaticFS))
	})

	apiRouter := r.Group("/api/bookmarks")
	{
		apiRouter.GET("", app.bookmarkController.GetAll)
		apiRouter.GET("/:id", app.bookmarkController.GetById)
		apiRouter.POST("", app.bookmarkController.Create)
		apiRouter.PUT("/:id", app.bookmarkController.Update)
		apiRouter.DELETE("/:id", app.bookmarkController.Delete)
	}

	return r
}

func (app *App) rootRouteHandler(c *gin.Context) {
	tmpl, err := template.ParseFS(assets.Templates, "templates/index.html")
	if err != nil {
		log.Fatalf("error loading static assets: %v", err)
	}
	tmpl.Execute(c.Writer, nil)
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
