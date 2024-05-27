package routers

import (
	_ "github.com/elef-git/chat_tool_golang/docs"
	"github.com/elef-git/chat_tool_golang/internal/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"time"
)

type middleware interface {
	RequireAuth() gin.HandlerFunc
}

type config interface {
	GetCorsAllowOrigins() []string
}

func InitRouter(env string, userV1Handler *handler.UserV1Handler, wsV1Handler *handler.WsV1Handler, postV1Handler *handler.PostV1Handler, m middleware, c config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     c.GetCorsAllowOrigins(),
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	if env == "dev" {
		r.Use(gin.Logger())
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	rootV1 := r.Group("/v1")
	authV1 := rootV1.Group("/auth")
	{
		authV1.POST("/login", userV1Handler.Login)
		authV1.POST("/registration", userV1Handler.Registration)
	}

	apiV1 := rootV1.Group("/api")
	apiV1.Use(m.RequireAuth())
	{
		apiV1.GET("/channels", userV1Handler.GetChannels)
		apiV1.GET("/posts/:channelID", postV1Handler.GetPosts)
		apiV1.POST("/posts", postV1Handler.AddPost)
		apiV1.PUT("/statuses", userV1Handler.UpdateStatus)
	}
	rootV1.GET("/ws/", m.RequireAuth(), wsV1Handler.NewWsConnection)

	return r
}
