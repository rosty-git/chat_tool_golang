package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/elef-git/chat_tool_golang/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type config interface {
	GetJwtSecret() string
	GetAuthCookieName() string
	GetJwtTtl() time.Duration

	GetAuthCookieMaxAge() int
	GetAuthCookiePath() string
	GetAuthCookieDomain() string
	GetAuthCookieSecure() bool
	GetAuthCookieHttpOnly() bool
}

type userRepository interface {
	GetById(db *gorm.DB, string string) (*models.User, error)
}

type Middleware struct {
	userRepository userRepository
	config         config
	db             *gorm.DB
}

func NewMiddleware(userRepository userRepository, config config, db *gorm.DB) *Middleware {
	return &Middleware{
		userRepository: userRepository,
		config:         config,
		db:             db,
	}
}

func (m *Middleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the cookie off the request
		tokenString, err := c.Cookie(m.config.GetAuthCookieName())
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

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

			sub, ok := claims["sub"]
			if !ok {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			if (claims["exp"].(float64)-float64(time.Now().Unix()))/m.config.GetJwtTtl().Seconds() < 0.5 {
				newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub": sub,
					"exp": time.Now().Add(m.config.GetJwtTtl()).Unix(),
				})

				// Sign and get the complete encoded token as a string using the secret
				newTokenString, err := newToken.SignedString([]byte(m.config.GetJwtSecret()))
				if err != nil {
					slog.Error("RequireAuth", "err", err)

					c.Error(err)
				}

				c.SetSameSite(http.SameSiteLaxMode)
				c.SetCookie(
					m.config.GetAuthCookieName(),
					newTokenString,
					m.config.GetAuthCookieMaxAge(),
					m.config.GetAuthCookiePath(),
					m.config.GetAuthCookieDomain(),
					m.config.GetAuthCookieSecure(),
					m.config.GetAuthCookieHttpOnly(),
				)
			}

			// Find the user with token Subject
			user, err := m.userRepository.GetById(m.db, sub.(string))
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			if user.ID == "" {
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
