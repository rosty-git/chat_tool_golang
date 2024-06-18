package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type GinConfig struct {
	Port string
}

type MySqlConfig struct {
	User   string
	Pass   string
	Host   string
	Port   string
	DbName string
}

type JwtConfig struct {
	Secret string
	Ttl    int
}

type CorsConfig struct {
	AllowOrigins []string
}

type AuthCookieConfig struct {
	Name          string
	MaxAgeSeconds int
	Path          string
	Domain        string
	Secure        bool
	HttpOnly      bool
}

type GormConfig struct {
	Debug bool
}

type AwsConfig struct {
	Region   string
	S3Bucket string
}

type Config struct {
	Env        string
	Gin        GinConfig
	MySql      MySqlConfig
	Jwt        JwtConfig
	Cors       CorsConfig
	AuthCookie AuthCookieConfig
	Gorm       GormConfig
	Aws        AwsConfig
}

func NewConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "prod"
	}

	ginPort := os.Getenv("GIN_PORT")
	if ginPort == "" {
		ginPort = "8080"
	}

	mySqlUser := os.Getenv("MYSQL_USER")
	if mySqlUser == "" {
		panic("MYSQL_USER env variable not set")
	}

	mySqlPass := os.Getenv("MYSQL_PASS")
	if mySqlPass == "" {
		panic("MYSQL_PASS env variable not set")
	}

	mySqlHost := os.Getenv("MYSQL_HOST")
	if mySqlHost == "" {
		mySqlHost = "localhost"
	}

	mySqlPort := os.Getenv("MYSQL_PORT")
	if mySqlPort == "" {
		mySqlPort = "3306"
	}

	mySqlDbName := os.Getenv("MYSQL_DB_NAME")
	if mySqlDbName == "" {
		panic("MYSQL_DB_NAME env variable not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET env variable not set")
	}

	jwtTtlSeconds := os.Getenv("JWT_TTL_SECONDS")
	if jwtTtlSeconds == "" {
		panic("JWT_TTL_SECONDS env variable not set")
	}
	jwtTtl, err := strconv.Atoi(jwtTtlSeconds)
	if err != nil {
		panic(fmt.Sprintf("JWT_TTL_SECONDS env variable not valid: %s", jwtTtlSeconds))
	}

	corsAllowOrigins := os.Getenv("CORS_ALLOW_ORIGINS")
	var corsAllowOriginsSlice []string
	if corsAllowOrigins == "" {
		corsAllowOriginsSlice = []string{"*"}
		slog.Warn("CORS_ALLOW_ORIGINS env variable has default value '*'")
	} else {
		corsAllowOriginsSlice = strings.Split(corsAllowOrigins, ",")
	}

	AuthCookieName := os.Getenv("AUTH_COOKIE_NAME")
	if AuthCookieName == "" {
		AuthCookieName = "Authorization"
	}

	AuthCookieMaxAgeSeconds := os.Getenv("AUTH_COOKIE_MAX_AGE_SECONDS")
	if AuthCookieMaxAgeSeconds == "" {
		AuthCookieMaxAgeSeconds = "2592000" // 30 days
	}
	AuthCookieMaxAge, err := strconv.Atoi(AuthCookieMaxAgeSeconds)
	if err != nil {
		panic(fmt.Sprintf("AUTH_COOKIE_MAX_AGE_SECONDS env variable not valid: %s", AuthCookieMaxAge))
	}

	AuthCookiePath := os.Getenv("AUTH_COOKIE_PATH")
	if AuthCookiePath == "" {
		AuthCookiePath = "/"
	}

	AuthCookieDomain := os.Getenv("AUTH_COOKIE_DOMAIN")
	if AuthCookieDomain == "" {
		AuthCookieDomain = "localhost"
	}

	AuthCookieSecure := os.Getenv("AUTH_COOKIE_SECURE")
	AuthCookieSecureBool, err := strconv.ParseBool(AuthCookieSecure)
	if err != nil {
		panic("AUTH_COOKIE_SECURE env variable not valid")
	}

	AuthCookieHttpOnly := os.Getenv("AUTH_COOKIE_HTTPONLY")
	AuthCookieHttpOnlyBool, err := strconv.ParseBool(AuthCookieHttpOnly)
	if err != nil {
		panic("AUTH_COOKIE_HTTPONLY env variable not valid")
	}

	GormDebugBool := false
	GormDebug := os.Getenv("GORM_DEBUG")
	if GormDebug != "" {
		GormDebugBool, err = strconv.ParseBool(GormDebug)
		if err != nil {
			panic("GORM_DEBUG env variable not valid")
		}
	}

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		panic("AWS_REGION env variable not set")
	}

	awsS3Bucket := os.Getenv("AWS_S3_BUCKET")
	if awsS3Bucket == "" {
		panic("AWS_S3_BUCKET env variable not set")
	}

	return &Config{
		Env: env,
		Gin: GinConfig{
			Port: ginPort,
		},
		MySql: MySqlConfig{
			User:   mySqlUser,
			Pass:   mySqlPass,
			Host:   mySqlHost,
			Port:   mySqlPort,
			DbName: mySqlDbName,
		},
		Jwt: JwtConfig{
			Secret: jwtSecret,
			Ttl:    jwtTtl,
		},
		Cors: CorsConfig{
			AllowOrigins: corsAllowOriginsSlice,
		},
		AuthCookie: AuthCookieConfig{
			Name:          AuthCookieName,
			MaxAgeSeconds: AuthCookieMaxAge,
			Path:          AuthCookiePath,
			Domain:        AuthCookieDomain,
			Secure:        AuthCookieSecureBool,
			HttpOnly:      AuthCookieHttpOnlyBool,
		},
		Gorm: GormConfig{
			Debug: GormDebugBool,
		},
		Aws: AwsConfig{
			Region:   awsRegion,
			S3Bucket: awsS3Bucket,
		},
	}
}

func (msc *MySqlConfig) ToDsnString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", msc.User, msc.Pass, msc.Host, msc.Port, msc.DbName)
}

func (c *Config) GetEnv() string {
	return c.Env
}

func (c *Config) GetJwtSecret() string {
	return c.Jwt.Secret
}

func (c *Config) GetJwtTtl() time.Duration {
	return time.Duration(c.Jwt.Ttl) * 1_000_000_000
}

func (c *Config) GetCorsAllowOrigins() []string {
	return c.Cors.AllowOrigins
}

func (c *Config) GetAuthCookieName() string {
	return c.AuthCookie.Name
}

func (c *Config) GetAuthCookiePath() string {
	return c.AuthCookie.Path
}

func (c *Config) GetAuthCookieDomain() string {
	return c.AuthCookie.Domain
}

func (c *Config) GetAuthCookieMaxAge() int {
	return c.AuthCookie.MaxAgeSeconds
}

func (c *Config) GetAuthCookieSecure() bool {
	return c.AuthCookie.Secure
}

func (c *Config) GetAuthCookieHttpOnly() bool {
	return c.AuthCookie.HttpOnly
}

func (c *Config) GetDsn() string {
	return c.MySql.ToDsnString()
}

func (c *Config) GetGormDebug() bool {
	return c.Gorm.Debug
}

func (c *Config) GetAwsRegion() string {
	return c.Aws.Region
}

func (c *Config) GetAwsS3Bucket() string {
	return c.Aws.S3Bucket
}
