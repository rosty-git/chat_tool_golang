package routers

import (
	"time"

	"github.com/elef-git/chat_tool_golang/internal/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type middleware interface {
	RequireAuth() gin.HandlerFunc
}

type config interface {
	GetCorsAllowOrigins() []string
}

func InitRouter(env string, userV1Handler *handler.UserV1Handler, wsV1Handler *handler.WsV1Handler, postV1Handler *handler.PostV1Handler, fileV1Handler *handler.FileV1Handler, m middleware, c config) *gin.Engine {
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
		apiV1.GET("/channels/:channelID/members", userV1Handler.GetChannelMembers)
		apiV1.PUT("/channels/:channelID/markasread", userV1Handler.MarkChannelAsRead)
		apiV1.GET("/channels/:channelID/unread", userV1Handler.GetUnreadCount)
		apiV1.GET("/channels/search/:text", userV1Handler.SearchChannels)
		apiV1.GET("/posts/:channelID", postV1Handler.GetPosts)
		apiV1.GET("/posts/search/:text", postV1Handler.SearchPosts)
		apiV1.POST("/posts", postV1Handler.AddPost)
		apiV1.PUT("/statuses", userV1Handler.UpdateStatus)
		apiV1.GET("/statuses/:userID", userV1Handler.GetStatus)
		apiV1.GET("/users/iam", userV1Handler.GetUser)
		apiV1.POST("/files", fileV1Handler.Create)
		apiV1.PUT("/files/:fileID/:s3Key", fileV1Handler.SetS3Key)
		apiV1.DELETE("/files/:fileID", fileV1Handler.Delete)
		apiV1.POST("/files/get-presigned-url/:key", fileV1Handler.GetPresignUrl)
	}
	rootV1.GET("/ws/", m.RequireAuth(), wsV1Handler.NewWsConnection)

	return r
}
