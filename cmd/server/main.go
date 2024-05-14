package main

import (
	"fmt"
	"github.com/elef-git/chat_tool_golang/internal/config"
	"github.com/elef-git/chat_tool_golang/internal/database"
	"github.com/elef-git/chat_tool_golang/internal/handler"
	"github.com/elef-git/chat_tool_golang/internal/middleware"
	"github.com/elef-git/chat_tool_golang/internal/repositories/user"
	"github.com/elef-git/chat_tool_golang/internal/routers"
	"github.com/elef-git/chat_tool_golang/internal/services/user"
	userusecase "github.com/elef-git/chat_tool_golang/internal/usecase/user"
	"github.com/elef-git/chat_tool_golang/pkg/logger"
	"log/slog"
	"net/http"
	"time"
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

	logger.InitLogger(c.Env)

	db, closer, err := database.New(c.MySql.ToDsnString(), c.Env)
	if err != nil {
		slog.Error("Failed to connect to database")
	}
	defer closer()

	err = database.Initialize(db)
	if err != nil {
		slog.Error("Failed to initialize database")
	}

	userRepo := userrepository.NewUsersRepository(db)
	userService := userservice.NewService(userRepo, c)
	userUseCase := userusecase.NewUseCase(userService)
	userV1Handler := handler.NewUserV1Handler(userUseCase, c)

	m := middleware.NewMiddleware(userRepo, c)
	router := routers.InitRouter(c.Env, userV1Handler, m, c)

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
