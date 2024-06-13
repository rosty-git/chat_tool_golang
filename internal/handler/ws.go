package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/elef-git/chat_tool_golang/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsMessage struct {
	ToUsersIDs []string
	Action     string
	Payload    interface{}
}

type WsV1Handler struct {
	config           config
	channel          chan WsMessage
	broadcastChannel chan WsMessage
	connMap          map[string]map[*websocket.Conn]bool
}

func NewWsV1Handler(config config, channel chan WsMessage, broadcastChannel chan WsMessage) *WsV1Handler {
	return &WsV1Handler{
		config:           config,
		channel:          channel,
		broadcastChannel: broadcastChannel,
		connMap:          make(map[string]map[*websocket.Conn]bool),
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

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		slog.Info("pong")
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	if _, ok := wh.connMap[user.ID]; !ok {
		wh.connMap[user.ID] = make(map[*websocket.Conn]bool)
	}

	wh.connMap[user.ID][conn] = true

	for {
		select {
		case msg := <-wh.channel:
			slog.Info("received", "msg", msg)

			for _, userId := range msg.ToUsersIDs {
				userConnections, ok := wh.connMap[userId]
				if ok {
					bytes, err := json.Marshal(msg)
					if err != nil {
						slog.Error("Marshal", "err", err)
					}

					for userConn, _ := range userConnections {
						err = userConn.WriteMessage(websocket.TextMessage, bytes)
						if err != nil {
							slog.Error("WriteMessage:", "err", err.Error())

							userConn.Close()

							delete(userConnections, userConn)
						}
					}
				} else {
					slog.Info("userConn not found")
				}
			}
		case bcMsg := <-wh.broadcastChannel:
			slog.Info("broadcast msg", "bcMsg", bcMsg)

			for userID, connections := range wh.connMap {
				slog.Info("WS Handler", "userID:", userID, "connections:", connections)

				bytes, err := json.Marshal(bcMsg)
				if err != nil {
					slog.Error("Marshal", "err", err)
				}

				for userConn, _ := range connections {
					err = userConn.WriteMessage(websocket.TextMessage, bytes)
					if err != nil {
						slog.Error("WriteMessage:", "err", err.Error())

						userConn.Close()

						delete(connections, userConn)
					}
				}
			}
		}
	}
}
