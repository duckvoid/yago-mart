package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

const (
	ansiReset          = "\033[0m"
	ansiBrightRed      = "\033[91m"
	ansiBrightGreen    = "\033[92m"
	ansiBrightYellow   = "\033[93m"
	ansiBrightRedFaint = "\033[91;2m"
)

func New(level string) *slog.Logger {
	lvl := &slog.LevelVar{}
	switch level {
	case "debug":
		lvl.Set(slog.LevelDebug)
	case "info":
		lvl.Set(slog.LevelInfo)
	case "warn":
		lvl.Set(slog.LevelWarn)
	case "error":
		lvl.Set(slog.LevelError)
	default:
		lvl.Set(slog.LevelInfo)
	}

	return slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level:      lvl,
			TimeFormat: time.DateTime,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.LevelKey {
					level := a.Value.Any().(slog.Level)
					switch {
					case level == slog.LevelError:
						a.Value = slog.StringValue(ansiBrightRed + "ERROR" + ansiReset)
					case level == slog.LevelWarn:
						a.Value = slog.StringValue(ansiBrightYellow + "WARN" + ansiReset)
					case level == slog.LevelInfo:
						a.Value = slog.StringValue(ansiBrightGreen + "INFO" + ansiReset)
					case level == slog.LevelDebug:
						a.Value = slog.StringValue(ansiBrightRedFaint + "DEBUG" + ansiReset)
					default:
						a.Value = slog.StringValue("UNKNOWN")
					}
				}

				return a
			},
		}),
	)
}
