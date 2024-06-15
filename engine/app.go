package engine

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"yc-backend/common"
	"yc-backend/internals"
	"yc-backend/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	Config  *common.Config
	DB      *mongo.Client
	Logger  internals.Logger
	Context context.Context
	server  *http.Server

	mux  *gin.Engine
	wg   sync.WaitGroup
	quit chan os.Signal
}

func (srv *Application) Setup() *Application {
	r := gin.New()

	r.Use(common.AddRequestIDMiddleware())
	r.Use(common.AddLoggerMiddleware(srv.Logger))
	r.Use(common.AddConfigMiddleware(srv.Config))
	r.Use(common.AddReposToMiddleware(srv.DB))
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(func(ctx *gin.Context) {
		cfg := common.GetConfigFromCtx(ctx)
		if cfg == nil {
			panic(errors.New("allowed cross-origin not set"))
		}
		cors.New(cors.Config{
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
			AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			ExposeHeaders:    []string{"Content-Length"},
			AllowOriginFunc: func(origin string) bool {
				return lo.Contains(cfg.AllowedCorsOrigin, origin)
			},
			MaxAge: 30 * time.Second,
		})(ctx)
		ctx.Next()
	})

	r.GET("/status-check", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, utils.SuccessResponse("alive🫵", gin.H{"status": "ok"}))
	})

	srv.mux = r
	srv.wg = sync.WaitGroup{}
	srv.server = &http.Server{
		Addr:              srv.Config.ServerPort,
		Handler:           srv.mux,
		ReadTimeout:       15 * time.Millisecond,
		WriteTimeout:      30 * time.Millisecond,
		ReadHeaderTimeout: 15 * time.Millisecond,
	}

	srv.quit = make(chan os.Signal, 1)
	signal.Notify(srv.quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	return srv
}

func (srv *Application) GracefulShutdown() {
	go func(quit chan os.Signal, dbm *mongo.Client) {
		<-quit
		shutdownCtx, shutdownCancelFunc := context.WithTimeout(srv.Context, 5*time.Second)
		go func() {
			defer shutdownCancelFunc()
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := srv.server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		err = dbm.Disconnect(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal("graceful shutdown... forcing exit.")
	}(srv.quit, srv.DB)
}

func (srv *Application) ListenAndServe() error {
	if err := srv.mux.Run(srv.Config.ServerPort); err != nil {
		return err
	}
	<-srv.Context.Done()
	return nil
}
