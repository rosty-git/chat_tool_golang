package handler

import (
	"encoding/json"
	"fmt"
	"github.com/elef-git/chat_tool_golang/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"slices"
)

type WsMessage struct {
	ToUsersIDs []string
	Action     string
	Payload    interface{}
}

type WsV1Handler struct {
	config  config
	channel chan WsMessage
	connMap map[string]*websocket.Conn
}

func NewWsV1Handler(config config, channel chan WsMessage) *WsV1Handler {
	return &WsV1Handler{
		config:  config,
		channel: channel,
		connMap: make(map[string]*websocket.Conn),
	}
}

func (wh *WsV1Handler) NewWsConnection(c *gin.Context) {
	slog.Info("NewWsConnection")

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return slices.Contains(wh.config.GetCorsAllowOrigins(), r.Header.Get("Origin"))
		},
	}

	authUser, ok := c.Get("user")
	if !ok {
		slog.Error("user not found")

		c.JSON(http.StatusInternalServerError, gin.H{})
	}

	user := authUser.(*models.User)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("conn", "err", err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
	}
	defer conn.Close()

	wh.connMap[user.ID] = conn

	for {
		select {
		case msg := <-wh.channel:
			fmt.Println("received", msg)

			for _, userId := range msg.ToUsersIDs {
				userConn, ok := wh.connMap[userId]
				if ok {
					bytes, err := json.Marshal(msg)
					if err != nil {
						slog.Error("Marshal", "err", err)
					}

					err = userConn.WriteMessage(websocket.TextMessage, bytes)
					if err != nil {
						slog.Error("WriteMessage:", "err", err.Error())
					}
				} else {
					slog.Info("userConn not found")
				}
			}

		}
	}
}
