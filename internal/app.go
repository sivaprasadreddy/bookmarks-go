package bookmarks

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/sivaprasadreddy/bookmarks-go/assets"
	"github.com/sivaprasadreddy/bookmarks-go/internal/api"
	"github.com/sivaprasadreddy/bookmarks-go/internal/config"
	"github.com/sivaprasadreddy/bookmarks-go/internal/db"
	"github.com/sivaprasadreddy/bookmarks-go/internal/domain"
	"github.com/sivaprasadreddy/bookmarks-go/internal/logging"
)

type App struct {
	Router             *gin.Engine
	cfg                config.AppConfig
	logger             *logging.Logger
	db                 *pgx.Conn
	bookmarkController *api.BookmarkController
}

func NewApp(cfg config.AppConfig) *App {
	app := &App{cfg: cfg}
	app.init()
	return app
}

func (app *App) init() {
	app.logger = logging.NewLogger(app.cfg)
	app.db = db.GetDb(app.cfg, app.logger)

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
		apiRouter.GET("/:id", app.bookmarkController.FindByID)
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
	// Create a context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	port := fmt.Sprintf(":%d", app.cfg.ServerPort)
	srv := &http.Server{
		Handler:        app.Router,
		Addr:           port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	app.logger.Infoln("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		app.logger.Fatal("Server forced to shutdown: ", err)
	}
	app.logger.Infoln("Server exiting")
}
