package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/elef-git/chat_tool_golang/internal/config"
	"github.com/elef-git/chat_tool_golang/internal/database"
	"github.com/elef-git/chat_tool_golang/internal/handler"
	"github.com/elef-git/chat_tool_golang/internal/middleware"
	channelrepository "github.com/elef-git/chat_tool_golang/internal/repositories/channel"
	filerepository "github.com/elef-git/chat_tool_golang/internal/repositories/file"
	postrepository "github.com/elef-git/chat_tool_golang/internal/repositories/post"
	"github.com/elef-git/chat_tool_golang/internal/repositories/user"
	"github.com/elef-git/chat_tool_golang/internal/routers"
	channelservice "github.com/elef-git/chat_tool_golang/internal/services/channel"
	fileservice "github.com/elef-git/chat_tool_golang/internal/services/file"
	postservice "github.com/elef-git/chat_tool_golang/internal/services/post"
	"github.com/elef-git/chat_tool_golang/internal/services/user"
	"github.com/elef-git/chat_tool_golang/internal/storages/s3"
	channelusecase "github.com/elef-git/chat_tool_golang/internal/usecase/channel"
	fileusecase "github.com/elef-git/chat_tool_golang/internal/usecase/file"
	postusecase "github.com/elef-git/chat_tool_golang/internal/usecase/post"
	userusecase "github.com/elef-git/chat_tool_golang/internal/usecase/user"
	"github.com/elef-git/chat_tool_golang/pkg/logger"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	c := config.NewConfig()

	logger.InitLogger(c)

	db, closer, err := database.New(c)
	if err != nil {
		slog.Error("Failed to connect to database")
	}
	defer closer()

	err = database.Initialize(db)
	if err != nil {
		slog.Error("Failed to initialize database")
	}

	wsChan := make(chan handler.WsMessage, 100)
	wsBroadcastChan := make(chan handler.WsMessage, 100)

	userRepo := userrepository.NewRepository()
	channelRepo := channelrepository.NewRepository()
	postRepo := postrepository.NewRepository()
	fileRepo := filerepository.NewRepository()

	userService := userservice.NewService(userRepo, c)
	channelService := channelservice.NewService(channelRepo)
	postService := postservice.NewService(postRepo, wsChan)
	fileService := fileservice.NewService(fileRepo)

	s3Storage := s3.NewStorage(c.GetAwsS3Bucket())

	userUseCase := userusecase.NewUseCase(db, userService, channelService, wsBroadcastChan)
	postUseCase := postusecase.NewUseCase(db, postService, channelService, fileService)
	channelUseCase := channelusecase.NewUseCase(db, channelService)
	fileUseCase := fileusecase.NewUseCase(db, fileService, s3Storage)

	go userUseCase.StatusesWatchdog()

	userV1Handler := handler.NewUserV1Handler(c, userUseCase, channelUseCase)
	postV1Handler := handler.NewPostV1Handler(postUseCase)
	wsV1Handler := handler.NewWsV1Handler(c, wsChan, wsBroadcastChan)
	fileV1Handler := handler.NewFileV1Handler(fileUseCase)

	m := middleware.NewMiddleware(userRepo, c, db)
	router := routers.InitRouter(c.Env, userV1Handler, wsV1Handler, postV1Handler, fileV1Handler, m, c)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", c.Gin.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	slog.Info(fmt.Sprintf("Start gin on port %s", c.Gin.Port))

	err = s.ListenAndServe()
	if err != nil {
		slog.Error("ListenAndServe", "err", err)
	}
}
