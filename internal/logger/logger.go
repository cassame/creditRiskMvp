package logger

import (
	"log/slog"
	"os"
)

var Lg *slog.Logger

func InitLogger(env string) {
	var handler slog.Handler
	if env == "local" {
		//local
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {
		//prod JSON
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	Lg = slog.New(handler)
	slog.SetDefault(Lg)
}
