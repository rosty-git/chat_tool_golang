package routers

import (
	"github.com/elef-git/chat_tool_golang/internal/handler"
	"github.com/elef-git/chat_tool_golang/internal/routers/api/v1"
	"github.com/gin-gonic/gin"

	_ "github.com/elef-git/chat_tool_golang/docs"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func InitRouter(env string, userV1Handler *handler.UserV1Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	if env == "dev" {
		r.Use(gin.Logger())
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	rootV1 := r.Group("/v1")
	authV1 := rootV1.Group("/auth")
	{
		authV1.GET("/login", userV1Handler.Login)
		authV1.POST("/registration", userV1Handler.Registration)
	}

	apiV1 := rootV1.Group("/api")
	//apiv1.Use(jwt.JWT())
	{
		apiV1.GET("/messages", v1.GetMessages)
	}

	return r
}
