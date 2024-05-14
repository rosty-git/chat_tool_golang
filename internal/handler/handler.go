package handler

type config interface {
	GetEnv() string
	GetAuthCookieName() string
	GetAuthCookieMaxAge() int
	GetAuthCookiePath() string
	GetAuthCookieDomain() string
	GetAuthCookieSecure() bool
	GetAuthCookieHttpOnly() bool
	GetCorsAllowOrigins() []string
}
