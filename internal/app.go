package bookmarks

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/sivaprasadreddy/bookmarks-go/assets"
	"github.com/sivaprasadreddy/bookmarks-go/internal/api"
	"github.com/sivaprasadreddy/bookmarks-go/internal/config"
	"github.com/sivaprasadreddy/bookmarks-go/internal/db"
	"github.com/sivaprasadreddy/bookmarks-go/internal/domain"
	"github.com/sivaprasadreddy/bookmarks-go/internal/logging"
	"html/template"
	"net/http"
	"path"
	"time"
)

type App struct {
	Router             *gin.Engine
	cfg                config.AppConfig
	logger             *logging.Logger
	db                 *pgx.Conn
	bookmarkController *api.BookmarkController
}

func NewApp(config config.AppConfig) *App {
	app := &App{}
	app.init(config)
	return app
}

func (app *App) init(config config.AppConfig) {
	app.logger = logging.NewLogger(config)
	app.db = db.GetDb(config, app.logger)

	bookmarksRepo := domain.NewBookmarkRepo(app.db, app.logger)
	app.bookmarkController = api.NewBookmarkController(bookmarksRepo, app.logger)

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
		apiRouter.GET("", app.bookmarkController.FindAll)
		apiRouter.GET("/:id", app.bookmarkController.FindById)
		apiRouter.POST("", app.bookmarkController.Create)
		apiRouter.PUT("/:id", app.bookmarkController.Update)
		apiRouter.DELETE("/:id", app.bookmarkController.Delete)
	}

	return r
}

func (app *App) rootRouteHandler(c *gin.Context) {
	tmpl, err := template.ParseFS(assets.Templates, "templates/index.html")
	if err != nil {
		app.logger.Fatalf("error loading static assets: %v", err)
	}
	err = tmpl.Execute(c.Writer, nil)
	if err != nil {
		app.logger.Fatalf("error rendering index.html: %v", err)
	}
}

func (app *App) Run() {
	port := fmt.Sprintf(":%d", app.cfg.ServerPort)
	srv := &http.Server{
		Handler:        app.Router,
		Addr:           port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	app.logger.Infof("listening on port %d", app.cfg.ServerPort)
	if err := srv.ListenAndServe(); err != nil {
		app.logger.Fatal(err)
	}
}
