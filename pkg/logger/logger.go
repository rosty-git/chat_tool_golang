package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type Config interface {
	GetEnv() string
}

func InitLogger(c Config) {
	if c.GetEnv() == "dev" {
		slog.SetDefault(slog.New(
			tint.NewHandler(os.Stderr, &tint.Options{
				Level:     slog.LevelDebug,
				AddSource: true,
			}),
		))
	} else {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	}
}
