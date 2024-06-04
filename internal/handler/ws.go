package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"slices"

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
	connMap          map[string][]*websocket.Conn
}

func NewWsV1Handler(config config, channel chan WsMessage, broadcastChannel chan WsMessage) *WsV1Handler {
	return &WsV1Handler{
		config:           config,
		channel:          channel,
		broadcastChannel: broadcastChannel,
		connMap:          make(map[string][]*websocket.Conn),
	}
}

func RemoveConnection(s []*websocket.Conn, index int) []*websocket.Conn {
	ret := make([]*websocket.Conn, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
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

	connectionsSlice := wh.connMap[user.ID]

	connectionsSlice = append(connectionsSlice, conn)

	wh.connMap[user.ID] = connectionsSlice

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

					for i, userConn := range userConnections {
						err = userConn.WriteMessage(websocket.TextMessage, bytes)
						if err != nil {
							slog.Error("WriteMessage:", "err", err.Error())

							userConn.Close()

							wh.connMap[userId] = RemoveConnection(userConnections, i)
						}
					}

					if len(wh.connMap[userId]) == 0 {
						delete(wh.connMap, userId)
					}
				} else {
					slog.Info("userConn not found")
				}
			}
		case bcMsg := <-wh.broadcastChannel:
			slog.Info("broadcast msg", "bcMsg", bcMsg)

			for userID, connections := range wh.connMap {
				fmt.Println("userID:", userID, "connections:", connections)

				bytes, err := json.Marshal(bcMsg)
				if err != nil {
					slog.Error("Marshal", "err", err)
				}

				for i, userConn := range connections {
					err = userConn.WriteMessage(websocket.TextMessage, bytes)
					if err != nil {
						slog.Error("WriteMessage:", "err", err.Error())

						userConn.Close()

						wh.connMap[userID] = RemoveConnection(connections, i)
					}
				}

				if len(wh.connMap[userID]) == 0 {
					delete(wh.connMap, userID)
				}
			}
		}
	}
}
