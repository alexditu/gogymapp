package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexditu/gogymapp/internal/logging"
	"github.com/alexditu/gogymapp/internal/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

type server struct {
	router *gin.Engine
	// dbConn    *mongo.Client
	// cfg       *config.Config
	ctx   context.Context
	setts *Settings

	localTZ string
}

func Run(binaryName string, setts *Settings) error {
	if setts == nil {
		fmt.Printf("Error: setts is nil")
		return fmt.Errorf("setts is nil")
	}

	if err := logging.InitStandardLogger(&setts.Log); err != nil {
		fmt.Printf("Failed to init logging: %s\n", err)
		return err
	}

	log.Infof("%s started (pid: %d)", binaryName, os.Getpid())

	// catch interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server, err := newServer(ctx, setts)
	defer server.stop()
	if err != nil {
		log.Errorf("init failed: %s", err)
		return err
	}

	server.setupRoutes()
	httpSrv := server.run()

	<-ctx.Done()

	// restore default behavior to interrupt signals
	stop()

	log.Info("Shutting down gracefully, press Ctrl+C again to force")

	// signal the server to graceful shutdown in max 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Warnf("Server forced to shutdown: %s", err)
		return err
	}

	log.Infof("%s finished (pid: %d)", binaryName, os.Getpid())

	return nil
}

func newServer(ctx context.Context, setts *Settings) (*server, error) {
	// use separate logging file for gin
	ginLogger, err := logging.NewLogger(&logging.Settings{LogToFile: true, FileName: "gin.log", Level: log.DebugLevel})
	if err != nil {
		return nil, err
	}

	// ginLogger.SetFormatter(&logrus.JSONFormatter{})

	router := gin.New()
	router.Use(ginlogrus.Logger(ginLogger), gin.Recovery())

	tz, err := utils.LocalTimeZone()
	if err != nil {
		log.Warnf("Failed to read timezone info: '%s'. Using UTC.", err)
		tz = "UTC"
	}

	s := server{router: router, ctx: ctx, localTZ: tz, setts: setts}

	return &s, nil
}

// server methods
func (s *server) stop() {
}

func (s *server) run() *http.Server {
	w := log.StandardLogger().Writer()

	// write newHttpServer in another file due to log naming conflict (logrus vs log)
	srv := newHttpServer(s, w)

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		err := srv.ListenAndServeTLS(s.setts.CertPath, s.setts.KeyPath)
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start https server: %s", err)
		}

		// close the pipe that http.Server is writing to after goroutine exits
		defer w.Close()
	}()

	return srv
}

func (s *server) setupRoutes() {
	s.router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// staticContentPath := s.getStaticContentPath()

	// s.router.Use(favicon.New(staticContentPath + "/images/logo_red_transparent.png"))

	// need by HTML Templating engine
	// s.router.LoadHTMLGlob(staticContentPath + "/*.html")

	// welcome page and static resources such as images
	s.router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/dashboard")
	})

	// s.router.Static("/static", staticContentPath)

	// REST API
	api := s.router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/login", s.login())
		}
	}

}

func (s *server) login() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
