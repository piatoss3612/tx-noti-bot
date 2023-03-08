package logger

import (
	"io"
	"path/filepath"

	"golang.org/x/exp/slog"
)

func SetStructuredLogger(name string, out io.Writer) {
	slog.SetDefault(NewStructuredLogger(name, out))
}

func NewStructuredLogger(name string, out io.Writer) *slog.Logger {
	opts := slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				a.Value = slog.StringValue(filepath.Base(a.Value.String()))
			}
			return a
		},
	}

	h := opts.NewJSONHandler(out)

	return slog.New(h).With("name", name)
}
