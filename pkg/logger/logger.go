package logger

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
)

func InitLogger(env string) {
	if env == "dev" {
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
