package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/elef-git/chat_tool_golang/internal/models"
	"github.com/gin-gonic/gin"
)

type PostV1Handler struct {
	postUseCase postUseCase
}

func NewPostV1Handler(postUseCase postUseCase) *PostV1Handler {
	return &PostV1Handler{
		postUseCase: postUseCase,
	}
}

func (uh *PostV1Handler) GetPosts(c *gin.Context) {
	slog.Info("UserV1Handler GetPosts", "channelID", c.Param("channelID"), "limit", c.Query("limit"))

	// TODO check if user is channel member

	var limitDefault int64 = 20
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		slog.Error("Error parsing limit", "err", err)
		limit = limitDefault
	}
	if limit > 100 {
		limit = limitDefault
	}

	//afterCreatedAt
	afterCreatedAtQuery := c.Query("afterCreatedAt")
	if afterCreatedAtQuery == "" {
		afterCreatedAtQuery = "1970-01-01T00:00:00.000Z"
	}

	afterCreatedAt, err := time.Parse(time.RFC3339, afterCreatedAtQuery)
	if err != nil {
		slog.Error("Error parsing afterCreatedAt", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid afterCreatedAt"})
		return
	}

	posts, err := uh.postUseCase.GetByChannelId(c.Param("channelID"), int(limit), afterCreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func (uh *PostV1Handler) AddPost(c *gin.Context) {
	slog.Info("PostV1Handler AddPost")

	type MessageForm struct {
		Channel string `json:"channel"`
		Message string `json:"message"`
	}

	var mf MessageForm
	err := c.BindJSON(&mf)
	if err != nil {
		slog.Error("BindJSON", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	slog.Info("BindJSON", "mf", mf)

	authUser, ok := c.Get("user")
	if !ok {
		slog.Error("user not found")
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	user := authUser.(*models.User)

	slog.Info("AddPost", "user", user)

	post, err := uh.postUseCase.Create(user.ID, mf.Channel, mf.Message)
	if err != nil {
		slog.Error("Create post", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}
