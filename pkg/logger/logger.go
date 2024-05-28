package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
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
