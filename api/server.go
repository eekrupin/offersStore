package api

import (
	"context"
	"github.com/eekrupin/offersStore/api/httpHandlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/eekrupin/offersStore/api/middlewares"
	"github.com/eekrupin/offersStore/config"
	"github.com/gin-gonic/gin"
	"github.com/semihalev/gin-stats"
	//"github.com/opentracing/opentracing-go"
	"github.com/gin-contrib/pprof"
)

func Run() {
	gin.DisableConsoleColor()
	if !config.Config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(stats.RequestStats())
	r.Use(middlewares.RequestMiddleware())
	r.Use(middlewares.ResponseMiddleware())
	r.Use(middlewares.RecoveryMiddleware())
	//r.Use(middlewares.TracerMiddleware(tracer))

	otherEP := r.Group("/")
	{
		otherEP.GET("/health", httpHandlers.Health)
		otherEP.POST("/health", httpHandlers.Health)

		otherEP.GET("/loadShareFile", httpHandlers.LoadShareFile)
		otherEP.GET("/loadUrl", httpHandlers.LoadUrlFile)
	}
	pprof.Register(r)

	srv := &http.Server{
		Addr:         config.Config.HTTPServer.Host + ":" + strconv.Itoa(int(config.Config.HTTPServer.InternalPort)),
		Handler:      r,
		ReadTimeout:  600 * time.Second,
		WriteTimeout: 600 * time.Second,
		//TLSConfig: tlsConfig,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")

}
