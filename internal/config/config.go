package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
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

type Config struct {
	Env   string
	Gin   GinConfig
	MySql MySqlConfig
}

func (msc *MySqlConfig) ToDsnString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", msc.User, msc.Pass, msc.Host, msc.Port, msc.DbName)
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
		slog.Error("MYSQL_USER env variable not set")
	}

	mySqlPass := os.Getenv("MYSQL_PASS")
	if mySqlPass == "" {
		slog.Error("MYSQL_PASS env variable not set")
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
		slog.Error("MYSQL_DB_NAME env variable not set")
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
	}
}
