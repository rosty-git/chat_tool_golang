package v1

import (
	"github.com/gin-gonic/gin"
	"log/slog"
)

// @BasePath /api/v1

// GetMessages PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /messages [get]
func GetMessages(c *gin.Context) {
	slog.Info("GetMessages")

	c.JSON(200, gin.H{
		"hello": "world!",
	})
}
