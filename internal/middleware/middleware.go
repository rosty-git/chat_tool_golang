package middleware

import (
	"fmt"
	"github.com/elef-git/chat_tool_golang/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
	"time"
)

type config interface {
	GetJwtSecret() string
}

type userRepository interface {
	GetById(string) (*models.User, error)
}

type Middleware struct {
	userRepository userRepository
	config         config
}

func NewMiddleware(userRepository userRepository, config config) *Middleware {
	return &Middleware{
		userRepository: userRepository,
		config:         config,
	}
}

func (m *Middleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("RequireAuth")
		//u, err := c.Cookie("U")
		//if err != nil {
		//	slog.Error("Get cookie err:", "err", err)
		//} else {
		//	slog.Info("Get cookie cookie:", "cookie", u)
		//}

		// Get the cookie off the request
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		slog.Info("tokenString:", tokenString)

		// Decode/validate it
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(m.config.GetJwtSecret()), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check the expiry date
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			id, ok := claims["sub"]
			if !ok {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			// Find the user with token Subject
			user, err := m.userRepository.GetById(id.(string))
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			if user.ID == 0 {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			// Attach the request
			c.Set("user", user)

			//Continue
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
