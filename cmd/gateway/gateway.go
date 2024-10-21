package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/reddtsai/reddservice/internal/global"
)

const (
	AccessControlMaxAge = 12 * time.Hour
)

type Gateway struct {
	Handler *gin.Engine
}

// @title Swagger Gateway
// @version 1.0
// @description This is a Gateway API.
// @host localhost
// @in header
// @name Authorization
// @tag.name auth
func NewGateway() *Gateway {
	engine := gin.New()
	engine.Use(loggerMiddleware())
	engine.Use(gin.Recovery())
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "authorization"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           AccessControlMaxAge,
	}))
	g := &Gateway{
		Handler: engine,
	}

	return g
}

func (g *Gateway) register(h *Handler) {
	v1 := g.Handler.Group("/v1")
	v1.POST("/sign-up", h.SignUp)
}

func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// TODO: prometheus

		c.Next()

		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}
		duration := time.Since(start)
		msg := fmt.Sprintf("[%s] %15s |%3d| %13s | %-7s %s", "gateway", c.ClientIP(), c.Writer.Status(), duration.String(), c.Request.Method, path)
		global.Logger.Info(msg, zap.Int("status", c.Writer.Status()), zap.String("method", c.Request.Method), zap.String("path", path), zap.Duration("duration", duration))
	}
}
