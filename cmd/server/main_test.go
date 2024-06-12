package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/elef-git/chat_tool_golang/internal/config"
	"github.com/elef-git/chat_tool_golang/internal/database"
	"github.com/elef-git/chat_tool_golang/internal/handler"
	"github.com/elef-git/chat_tool_golang/internal/middleware"
	channelrepository "github.com/elef-git/chat_tool_golang/internal/repositories/channel"
	postrepository "github.com/elef-git/chat_tool_golang/internal/repositories/post"
	userrepository "github.com/elef-git/chat_tool_golang/internal/repositories/user"
	"github.com/elef-git/chat_tool_golang/internal/routers"
	channelservice "github.com/elef-git/chat_tool_golang/internal/services/channel"
	postservice "github.com/elef-git/chat_tool_golang/internal/services/post"
	userservice "github.com/elef-git/chat_tool_golang/internal/services/user"
	channelusecase "github.com/elef-git/chat_tool_golang/internal/usecase/channel"
	postusecase "github.com/elef-git/chat_tool_golang/internal/usecase/post"
	userusecase "github.com/elef-git/chat_tool_golang/internal/usecase/user"
	"github.com/elef-git/chat_tool_golang/pkg/logger"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func TestApp(t *testing.T) {
	ctx := context.Background()

	os.Setenv("ENV", "dev")
	os.Setenv("MYSQL_USER", "root")
	os.Setenv("MYSQL_PASS", "pass")
	os.Setenv("MYSQL_DB_NAME", "chat")
	os.Setenv("JWT_SECRET", "jwt_56cr6t")
	os.Setenv("JWT_TTL_SECONDS", "604800")
	os.Setenv("AUTH_COOKIE_SECURE", "true")
	os.Setenv("AUTH_COOKIE_HTTPONLY", "true")

	c := config.NewConfig()
	logger.InitLogger(c)

	mysqlContainer, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:8.4"),
		mysql.WithDatabase(c.MySql.DbName),
		mysql.WithUsername(c.MySql.User),
		mysql.WithPassword(c.MySql.Pass),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	port, err := mysqlContainer.MappedPort(ctx, nat.Port(c.MySql.Port))
	if err != nil {
		panic(err)
	}

	c.MySql.Port = port.Port()

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

	userService := userservice.NewService(userRepo, c)
	channelService := channelservice.NewService(channelRepo)
	postService := postservice.NewService(postRepo, wsChan)

	userUseCase := userusecase.NewUseCase(db, userService, channelService, wsBroadcastChan)
	postUseCase := postusecase.NewUseCase(db, postService, channelService)
	channelUseCase := channelusecase.NewUseCase(db, channelService)

	userV1Handler := handler.NewUserV1Handler(c, userUseCase, channelUseCase)
	postV1Handler := handler.NewPostV1Handler(postUseCase)
	wsV1Handler := handler.NewWsV1Handler(c, wsChan, wsBroadcastChan)

	m := middleware.NewMiddleware(userRepo, c, db)
	router := routers.InitRouter(c.Env, userV1Handler, wsV1Handler, postV1Handler, m, c)

	type LoginForm struct {
		Email    string
		Password string
	}

	type Channel struct {
		Id   string `json:"id"`
		Type string `json:"type"`
		Name string `json:"name"`
	}

	type ChannelsResponse struct {
		Channels []Channel `json:"channels"`
	}

	type PostForm struct {
		Message   string
		ChannelId string
	}

	type ChannelUnreadResponse struct {
		Unread int `json:"unread"`
	}

	type Post struct {
		Id string `json:"id"`
	}

	type PostsResponse struct {
		Posts []Post `json:"posts"`
	}

	type User struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	type UserIamResponse struct {
		User User `json:"user"`
	}

	type StatusForm struct {
		Status string
		Manual bool
	}

	type Status struct {
		Status string `json:"status"`
		Manual bool   `json:"manual"`
	}

	type StatusResponse struct {
		Status Status `json:"status"`
	}

	var authCookie *http.Cookie
	var user1DirectChannel Channel
	var userId string

	t.Run("LoginUser3", func(t *testing.T) {
		company := LoginForm{
			Email:    "user3@gmail.com",
			Password: "password3",
		}
		jsonValue, err := json.Marshal(company)
		req, err := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		cookie := w.Result().Cookies()
		authCookie = cookie[0]

		require.NoError(t, err)
		require.Equal(t, 200, w.Code)
		require.Equal(t, "Authorization", cookie[0].Name)
	})

	t.Run("GetChannels", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/api/channels?channelType=D", nil)
		req.AddCookie(authCookie)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var channelsResponse ChannelsResponse
		err = json.Unmarshal([]byte(w.Body.String()), &channelsResponse)

		for _, channel := range channelsResponse.Channels {
			if channel.Name == "user1" {
				user1DirectChannel = channel
			}
		}

		require.NoError(t, err)
		require.Equal(t, 200, w.Code)
		require.Equal(t, 2, len(channelsResponse.Channels))
		require.Equal(t, "D", channelsResponse.Channels[0].Type)
		require.Equal(t, "user1", user1DirectChannel.Name)
	})

	t.Run("GetChannelMembers", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/api/channels/"+user1DirectChannel.Id+"/members", nil)
		req.AddCookie(authCookie)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		type ChannelMembersResponse struct {
			Members []string `json:"members"`
		}
		var channelMembersResponse ChannelMembersResponse

		err = json.Unmarshal([]byte(w.Body.String()), &channelMembersResponse)

		require.NoError(t, err)
		require.Equal(t, 200, w.Code)
		require.Equal(t, 1, len(channelMembersResponse.Members))
	})

	t.Run("LoginUser1", func(t *testing.T) {
		company := LoginForm{
			Email:    "user1@gmail.com",
			Password: "password1",
		}
		jsonValue, err := json.Marshal(company)
		req, err := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		cookieUser1 := w.Result().Cookies()

		require.NoError(t, err)
		require.Equal(t, 200, w.Code)
		require.Equal(t, "Authorization", cookieUser1[0].Name)

		t.Run("AddPost", func(t *testing.T) {
			post := PostForm{
				ChannelId: user1DirectChannel.Id,
				Message:   "HW",
			}
			postValue, err := json.Marshal(post)
			req, err := http.NewRequest("POST", "/v1/api/posts", bytes.NewBuffer(postValue))
			req.AddCookie(cookieUser1[0])
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			//slog.Info("", "", w.Body)

			require.NoError(t, err)
			require.Equal(t, 200, w.Code)
		})
	})

	t.Run("GetChannelUnread", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/api/channels/"+user1DirectChannel.Id+"/unread", nil)
		req.AddCookie(authCookie)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var channelUnreadResponse ChannelUnreadResponse
		err = json.Unmarshal([]byte(w.Body.String()), &channelUnreadResponse)

		require.NoError(t, err)
		require.Equal(t, 200, w.Code)
		require.Equal(t, 1, channelUnreadResponse.Unread)
	})

	t.Run("MarkChannelRead", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/v1/api/channels/"+user1DirectChannel.Id+"/markasread", nil)
		req.AddCookie(authCookie)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.NoError(t, err)
		require.Equal(t, 200, w.Code)
	})

	t.Run("GetChannelUnreadAgain", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/api/channels/"+user1DirectChannel.Id+"/unread", nil)
		req.AddCookie(authCookie)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var channelUnreadResponse ChannelUnreadResponse
		err = json.Unmarshal([]byte(w.Body.String()), &channelUnreadResponse)

		require.NoError(t, err)
		require.Equal(t, 200, w.Code)
		require.Equal(t, 0, channelUnreadResponse.Unread)
	})

	t.Run("GetPosts", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/api/posts/"+user1DirectChannel.Id, nil)
		req.AddCookie(authCookie)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var postsResponse PostsResponse
		err = json.Unmarshal([]byte(w.Body.String()), &postsResponse)

		require.NoError(t, err)
		require.Equal(t, 200, w.Code)
		require.Equal(t, 1, len(postsResponse.Posts))
	})

	t.Run("GetUserIam", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/api/users/iam", nil)
		req.AddCookie(authCookie)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var userIamResponse UserIamResponse
		err = json.Unmarshal([]byte(w.Body.String()), &userIamResponse)

		userId = userIamResponse.User.Id

		require.NoError(t, err)
		require.Equal(t, 200, w.Code)
		require.Equal(t, "user3", userIamResponse.User.Name)
	})

	t.Run("PutStatus", func(t *testing.T) {
		status := StatusForm{
			Status: "dnd",
			Manual: true,
		}
		jsonValue, err := json.Marshal(status)

		req, err := http.NewRequest("PUT", "/v1/api/statuses", bytes.NewBuffer(jsonValue))
		req.AddCookie(authCookie)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.NoError(t, err)
		require.Equal(t, 200, w.Code)
	})

	t.Run("GetStatus", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/api/statuses/"+userId, nil)
		req.AddCookie(authCookie)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var statusResponse StatusResponse
		err = json.Unmarshal([]byte(w.Body.String()), &statusResponse)

		require.NoError(t, err)
		require.Equal(t, 200, w.Code)
		require.Equal(t, "dnd", statusResponse.Status.Status)
		require.Equal(t, true, statusResponse.Status.Manual)
	})

	// Clean up the container
	defer func() {
		if err := mysqlContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()
}
