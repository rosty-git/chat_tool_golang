package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"slices"
	"time"
)

type WsV1Handler struct {
	config config
}

func NewWsV1Handler(config config) *WsV1Handler {
	return &WsV1Handler{
		config: config,
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

	user, ok := c.Get("user")
	if !ok {
		slog.Error("user not found")
		c.JSON(http.StatusInternalServerError, "")
	}

	slog.Info("NewWsConnection", "User", user)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("conn", "err", err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
	}
	defer conn.Close()

	slog.Info("After conn")

	i := 0
	for {
		slog.Info("In for")
		i++

		err := conn.WriteMessage(websocket.TextMessage, []byte(`{"hello":100}`))
		if err != nil {
			slog.Error("upgrade:", "err", err.Error())
		}
		time.Sleep(time.Second * 5)
	}
}
