package routers

import (
	_ "github.com/elef-git/chat_tool_golang/docs"
	"github.com/elef-git/chat_tool_golang/internal/handler"
	"github.com/elef-git/chat_tool_golang/internal/routers/api/v1"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type middleware interface {
	RequireAuth() gin.HandlerFunc
}

func InitRouter(env string, userV1Handler *handler.UserV1Handler, m middleware) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.Default())

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
	//apiV1.Use(m.RequireAuth())
	{
		apiV1.GET("/messages", m.RequireAuth(), v1.GetMessages)
		//apiV1.GET("/messages", v1.GetMessages)
	}

	return r
}
