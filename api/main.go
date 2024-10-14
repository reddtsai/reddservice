package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/reddtsai/reddservice/api/gateway"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	g := gin.Default()
	docs.SwaggerInfo.BasePath = "/v1"
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	srv := http.Server{
		Addr:    ":8080",
		Handler: g,
	}
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown: %v", err)
	}
}
