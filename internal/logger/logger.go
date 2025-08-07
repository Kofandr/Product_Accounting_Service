package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/Kofandr/Product_Accounting_Service/internal/middleware"
)

func New(level string) *slog.Logger {
	otps := &slog.HandlerOptions{}

	switch level {
	case "DEBUG":
		otps.Level = slog.LevelDebug
	case "WARN":
		otps.Level = slog.LevelWarn
	case "ERROR":
		otps.Level = slog.LevelError
	default:
		otps.Level = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, otps))
}

func MustLoggerFromCtx(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(middleware.CtxLoggerKey{}).(*slog.Logger); ok {
		return logger
	}

	return slog.Default()
}
