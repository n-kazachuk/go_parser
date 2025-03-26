package sl

import (
	"log/slog"

	"github.com/n-kazachuk/go_parser/internal/libs/helpers"
)

func WithTrace(log *slog.Logger) *slog.Logger {
	return log.With(slog.String("trace", helpers.GetFunctionName(helpers.DefaultDepth+1)))
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
